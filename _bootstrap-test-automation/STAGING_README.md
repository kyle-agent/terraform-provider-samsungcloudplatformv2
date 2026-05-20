# Staging — DELETE AFTER PUSH

이 디렉터리는 `kyle-agent/terraform-test-automation` repo로 옮길 작업물의 임시 보관소입니다. 본 provider repo의 일부가 아니며, push 완료 후 별도 commit으로 삭제됩니다.

## 안에 든 것

- `terraform-test-automation.bundle` — 모든 파일이 담긴 git bundle (단일 commit `5d82c51`). 가장 손쉬운 이전 경로.
- 나머지 파일·디렉터리 — bundle을 unzip한 동일 내용. bundle 이전이 안 될 때 폴더 단위 복사 fallback.

## 새 세션에서 가져가는 절차 (택1)

### 방법 1 — bundle 이용 (권장)

```bash
# 1. provider repo를 임시로 받아 bundle 파일만 추출
git clone --depth 1 \
  --branch claude/fix-resource-deletion-bug-Vz59x \
  https://github.com/kyle-agent/terraform-provider-samsungcloudplatformv2.git /tmp/_p
cp /tmp/_p/_bootstrap-test-automation/terraform-test-automation.bundle /tmp/

# 2. 현재 작업 repo(빈 terraform-test-automation)로 이동
cd <terraform-test-automation 작업 디렉터리>

# 3. bundle 검증 + fetch + push
git bundle verify /tmp/terraform-test-automation.bundle
git fetch /tmp/terraform-test-automation.bundle main:main
git push -u origin main
```

### 방법 2 — 디렉터리 복사 (bundle이 깨졌거나 fetch가 거부될 때)

```bash
# 1. provider repo의 _bootstrap-test-automation/ 폴더만 받기
git clone --depth 1 \
  --branch claude/fix-resource-deletion-bug-Vz59x \
  --filter=blob:none --sparse \
  https://github.com/kyle-agent/terraform-provider-samsungcloudplatformv2.git /tmp/_p
cd /tmp/_p
git sparse-checkout set _bootstrap-test-automation

# 2. 현재 작업 repo로 파일 복사 (단, .bundle은 제외)
cd <terraform-test-automation 작업 디렉터리>
rsync -a --exclude '*.bundle' --exclude 'README.md' /tmp/_p/_bootstrap-test-automation/ ./

# 3. (README.md는 별도로 — bundle README가 아닌 본래 README 사용)
git checkout HEAD -- README.md 2>/dev/null || true
# (작업 repo가 빈 상태라면 위 라인 무시)

# 4. 초기 commit & push
git add -A
git -c commit.gpgsign=false commit -m "Bootstrap regression test automation system"
git push -u origin main
```

> 주의: 방법 2의 README.md는 본 staging의 README가 아니라 `_bootstrap-test-automation/README.md`(원본 README) 입니다. 파일명이 같으므로 rsync 대상에서 본 staging README는 자동으로 덮어쓰여 정상입니다.

## 정리 (push 성공 후)

```bash
git rm -r _bootstrap-test-automation
git commit -m "Remove bootstrap staging — pushed to terraform-test-automation repo"
git push
```
