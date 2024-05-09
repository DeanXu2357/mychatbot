# mychatbot(暫)

## 目錄

- [目錄](#目錄)
- [開始](#開始)
  - [專案目的](#專案目的)
  - [安裝](#安裝)
- [待辦事項](#待辦事項)

## 開始

### 專案目的
個人助手整合

### 安裝

#### 環境需求

## 待辦事項
- [ ] 功能模組
  - [ ] Workout 功能模組開發
    - [ ] 資料表設計
    - [ ] 訓練組數紀錄
    - [ ] 訓練量報表
    - [ ] 訓練計畫相關功能
  - [ ] 記帳功能模組開發
    - [ ] beancount 串接
    - [ ] github 專案作為資料庫整合
  - [ ] 個人行事曆與 google calendar 整合
  - [ ] 團購分單功能
    - [ ] 功能設計
- [ ] 交互接口開發
  - [ ] Line bot
  - [ ] Discord bot
- [ ] llm 相關
  - [ ] langchaingo 研究
- [ ] observability
  - [ ] prometheus
  - [ ] grafana
- [ ] Robustness
  - [ ] command `server` 中，有任何服務失敗時觸發 ctx 的 cancel 事件，graceful shutdown server
  - [ ] 完善 log 機制
