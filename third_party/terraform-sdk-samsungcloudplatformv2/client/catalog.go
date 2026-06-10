package scpsdk

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"slices"
	"sort"
	"strconv"
	"sync"
	"time"
)

const DEPRECATED_SDK_MSG_PREFIX = "Supported Until - Warning"
const CHECK_DEPRECATED_DAYS = 60

type Endpoint struct {
	Region      string `json:"region"`
	ServiceType string `json:"service_type"`
	ServiceName string `json:"service_name"`
	URL         string `json:"url"`
}

type Catalog struct {
	AuthURL       string
	AccessKey     string
	SecretKey     string
	DefaultRegion string
}

var catalog []Endpoint

type VersionLink struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

type Version struct {
	ID        string        `json:"id"`
	Links     []VersionLink `json:"links"`
	NotBefore string        `json:"not_before"`
	Status    string        `json:"status"`
}

type VersionsResponse struct {
	Versions []Version `json:"versions"`
}

func containsVersion(versions []Version, target string) bool {
	for _, v := range versions {
		if v.ID == target {
			return true
		}
	}
	return false
}

func NewCatalog(authURL, accessKey, secretKey, defaultRegion string) *Catalog {
	return &Catalog{
		AuthURL:       authURL,
		AccessKey:     accessKey,
		SecretKey:     secretKey,
		DefaultRegion: defaultRegion,
	}
}

func (c *Catalog) GetEndpointList(serviceTypeMap map[string][]string, accountID string) ([]Endpoint, error) {
	if len(catalog) == 0 {
		//resultList := []string{}
		req, err := http.NewRequest("GET", c.AuthURL+"/endpoints", nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return nil, err
		}

		timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		signature := MakeHmacSignature("GET", c.AuthURL+"/endpoints", timestamp, c.AccessKey, c.SecretKey, accountID, ClientType)

		req.Header.Set(HeaderClientType, ClientType)
		req.Header.Set(HeaderLanguage, "en-US")
		req.Header.Set(HeaderAccountId, accountID)
		req.Header.Set(HeaderAccessKey, c.AccessKey)
		req.Header.Set(HeaderTimestamp, timestamp)
		req.Header.Set(HeaderSignature, signature)

		certPath := os.Getenv("SSL_CERT_FILE")
		var certPool *x509.CertPool

		if certPath == "" {
			certPool, err = x509.SystemCertPool()
			if err != nil {
				return nil, err
			}
		} else {
			crt, err := ioutil.ReadFile(certPath)
			if err != nil {
				return nil, err
			}
			certPool = x509.NewCertPool()
			certPool.AppendCertsFromPEM(crt)
		}

		httpClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{RootCAs: certPool},
				Proxy:           http.ProxyFromEnvironment,
			},
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return nil, err
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error get endpoints response parsing:", err)
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Error get endpoints response:", string(body))
			return nil, errors.New(string(body))
		}

		var response struct {
			Endpoints []Endpoint `json:"endpoints"`
		}
		err = json.Unmarshal(body, &response)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return nil, err
		}

		for _, endpoint := range response.Endpoints {
			catalog = append(catalog, endpoint)
		}
	}

	var candidates []Endpoint

	for serviceType, _ := range serviceTypeMap {
		for _, endpoint := range catalog {
			if endpoint.ServiceType != serviceType {
				continue
			}
			candidates = append(candidates, endpoint)
		}
	}
	//
	//if len(candidates) == 0 {
	//	return "", errors.New("no matching endpoint found")
	//} else if len(candidates) == 1 {
	//	return candidates[0].URL, nil
	//}
	//
	//originalCandidates := make([]Endpoint, len(candidates))
	//copy(originalCandidates, candidates)
	//candidates = []Endpoint{}
	//
	//if region == "" {
	//	region = c.DefaultRegion
	//
	//if region != "" {
	//	for _, endpoint := range originalCandidates {
	//		if endpoint.Region != region {
	//			continue
	//		}
	//		candidates = append(candidates, endpoint)
	//	}
	//}
	//
	//if len(candidates) == 0 {
	//	candidates = originalCandidates
	//}
	//
	//sort.Slice(candidates, func(i, j int) bool {
	//	return candidates[i].URL < originalCandidates[j].URL
	//})
	//
	//return candidates[0].URL, nil
	return candidates, nil
}

func (c *Catalog) GetEndpoint(serviceType, region, accountID string) (string, error) {
	if len(catalog) == 0 {
		req, err := http.NewRequest("GET", c.AuthURL+"/endpoints", nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return "", err
		}

		timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		signature := MakeHmacSignature("GET", c.AuthURL+"/endpoints", timestamp, c.AccessKey, c.SecretKey, accountID, ClientType)

		req.Header.Set(HeaderClientType, ClientType)
		req.Header.Set(HeaderLanguage, "en-US")
		req.Header.Set(HeaderAccountId, accountID)
		req.Header.Set(HeaderAccessKey, c.AccessKey)
		req.Header.Set(HeaderTimestamp, timestamp)
		req.Header.Set(HeaderSignature, signature)

		certPath := os.Getenv("SSL_CERT_FILE")
		var certPool *x509.CertPool

		if certPath == "" {
			certPool, err = x509.SystemCertPool()
			if err != nil {
				return "", err
			}
		} else {
			crt, err := ioutil.ReadFile(certPath)
			if err != nil {
				return "", err
			}
			certPool = x509.NewCertPool()
			certPool.AppendCertsFromPEM(crt)
		}

		httpClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{RootCAs: certPool},
				Proxy:           http.ProxyFromEnvironment,
			},
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return "", err
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error get endpoints response parsing:", err)
			return "", err
		}
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Error get endpoints response:", string(body))
			return "", errors.New(string(body))
		}

		var response struct {
			Endpoints []Endpoint `json:"endpoints"`
		}
		err = json.Unmarshal(body, &response)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return "", err
		}

		for _, endpoint := range response.Endpoints {
			catalog = append(catalog, endpoint)
		}
	}

	var candidates []Endpoint

	for _, endpoint := range catalog {
		if endpoint.ServiceType != serviceType {
			continue
		}
		candidates = append(candidates, endpoint)
	}

	if len(candidates) == 0 {
		return "", errors.New("no matching endpoint found")
	} else if len(candidates) == 1 {
		return candidates[0].URL, nil
	}

	originalCandidates := make([]Endpoint, len(candidates))
	copy(originalCandidates, candidates)
	candidates = []Endpoint{}

	if region == "" {
		region = c.DefaultRegion
	}

	if region != "" {
		for _, endpoint := range originalCandidates {
			if endpoint.Region != region {
				continue
			}
			candidates = append(candidates, endpoint)
		}
	}

	if len(candidates) == 0 {
		candidates = originalCandidates
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].URL < originalCandidates[j].URL
	})

	return candidates[0].URL, nil
}

func (c *Catalog) CheckVersion(basepath string, allowVersion []string, serviceName string) (bool, error) {
	req, err := http.NewRequest("GET", basepath, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false, err
	}

	certPath := os.Getenv("SSL_CERT_FILE")
	var certPool *x509.CertPool

	if certPath == "" {
		certPool, err = x509.SystemCertPool()
		if err != nil {
			return false, err
		}
	} else {
		crt, err := ioutil.ReadFile(certPath)
		if err != nil {
			return false, err
		}
		certPool = x509.NewCertPool()
		certPool.AppendCertsFromPEM(crt)
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: certPool},
			Proxy:           http.ProxyFromEnvironment,
		},
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return false, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error get version response parsing:", err)
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error get version response:", string(body))
		return false, errors.New(string(body))
	}

	var response VersionsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return false, err
	}

	for _, allowVersion := range allowVersion {
		if containsVersion(response.Versions, allowVersion) == false {
			return false, errors.New(serviceName + " " + allowVersion + " is an Unsupported version. " +
				"please update your provider.")
		}
	}

	return true, nil
}

func (c *Catalog) AsyncVersionCheck(serviceList map[string][]string, endpoint Endpoint, results chan<- string, wg *sync.WaitGroup, nowDate time.Time, remainDays int64, checkTimeout int64) {
	defer wg.Done()

	basepath := endpoint.URL
	serviceName := endpoint.ServiceName
	region := endpoint.Region

	req, err := http.NewRequest("GET", basepath, nil)
	if err != nil {
		results <- fmt.Sprintf("\t%s (%s) Status check failed - %v", serviceName, region, err)
		return
	}

	certPath := os.Getenv("SSL_CERT_FILE")
	var certPool *x509.CertPool

	if certPath == "" {
		certPool, err = x509.SystemCertPool()
		if err != nil {
			results <- fmt.Sprintf("\t%s (%s) Status check failed - %v", serviceName, region, err)
			return
		}
	} else {
		crt, err := ioutil.ReadFile(certPath)
		if err != nil {
			results <- fmt.Sprintf("\t%s (%s) Status check failed - %v", serviceName, region, err)
			return
		}
		certPool = x509.NewCertPool()
		certPool.AppendCertsFromPEM(crt)
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: certPool},
			Proxy:           http.ProxyFromEnvironment,
		},
		Timeout: time.Duration(checkTimeout) * time.Second,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		results <- fmt.Sprintf("\t%s (%s) Status check failed - %v", serviceName, region, err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		results <- fmt.Sprintf("\t%s (%s) Status check failed - %v", serviceName, region, err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		err := "Http response " + resp.Status
		results <- fmt.Sprintf("\t%s (%s) Status check failed - %v", serviceName, region, err)
		return
	}

	var response VersionsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		results <- fmt.Sprintf("\t%s (%s) Status check failed - %v", serviceName, region, err)
		return
	}
	versioncheckResult := ""
	for _, ver := range response.Versions {
		if ver.NotBefore != "" {
			allowServiceVersionList := serviceList[serviceName]

			if slices.Contains(allowServiceVersionList, ver.ID) {
				expireDate, err := time.Parse("2006-01-02", ver.NotBefore)
				if err != nil {
					results <- fmt.Sprintf("\t%s (%s) Status check failed - %v", serviceName, region, err)
					return
				}
				daysLeft := int64(expireDate.Sub(nowDate).Hours() / 24)

				if daysLeft <= remainDays {
					versioncheckResult += fmt.Sprintf("\t%s %s(%s) [%s], Until Supported date %s: %d days left\n", DEPRECATED_SDK_MSG_PREFIX, serviceName, region, ver.ID, ver.NotBefore, daysLeft)
				}
			}
		}
	}

	if versioncheckResult != "" {
		results <- versioncheckResult
	}
	return
}
