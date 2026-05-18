# 🚀 Point Hub

統合ポイント管理ツール (Unified Transit Point Aggregator)

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)
[![Status](https://img.shields.io/badge/Status-Development-orange?style=for-the-badge)](https://github.com/kikudesuyo/point-hub)

PASMOを対象とした首都圏の私鉄各社のポイントを、コマンドライン一つで集約・管理するためのツールです。

## ✨ 特徴 (Features)

- **一括取得**: 複数の鉄道事業者のポイント残高と有効期限を一度に取得。
- **モダンなWeb UI**: SvelteKit/Tailwind CSSを使用したプレミアムなダッシュボード。
- **高速な動作**: `localStorage` を活用したキャッシュ機構により、表示速度を最適化。
- **統一フォーマット**: 各社バラバラなデータ形式を、見やすい統一された形式で表示。

## 🚉 対応プロバイダー (Supported Providers)

| プロバイダー   | サービス名                                                  | 認証方式             | 備考                                     |
| :------------- | :---------------------------------------------------------- | :------------------- | :--------------------------------------- |
| **東京メトロ** | [メトポ (Metpo)](docs/tokyo_metro.md)                       | メール・パスワード   | 毎月11日付与, [詳細](docs/tokyo_metro.md) |
| **都営地下鉄** | [ToKoPo (トコポ)](docs/toei_metro.md)                       | 会員番号・パスワード | [詳細](docs/toei_metro.md)               |
| **東急電鉄**   | [東急ポイント](docs/tokyu.md)                               | セッショントークン   | [詳細](docs/tokyu.md)                    |
| **相模鉄道**   | [相鉄ポイント](docs/sotetsu.md)                             | メール・パスワード   | [詳細](docs/sotetsu.md)                  |
| **京浜急行**   | [京急ポイント](docs/keikyu.md)                               | 会員番号・パスワード | 毎月10日付与, [詳細](docs/keikyu.md)     |
| **小田急電鉄** | [小田急ポイント](docs/odakyu.md)                             | Auth0 Bearerトークン | [NEW]                                    |
| **東武鉄道**   | [東武ポイント](docs/tobu.md)                                 | セッションクッキー   | [NEW]                                    |

### 📅 ポイント付与スケジュール

- **東京メトロ (Metpo)**: 毎月11日
- **京浜急行**: 毎月10日

## 🛠 セットアップ (Setup)

### 1. リポジトリのクローン

```bash
git clone https://github.com/kikudesuyo/point-hub.git
cd point-hub
```

### 2. バックエンド (Go) の起動

```bash
cd api
go mod tidy
go run main.go # Port 8081で起動
```

### 3. フロントエンド (SvelteKit) の起動

```bash
cd frontend
npm install
npm run dev # Port 5173で起動
```

## 🚀 使い方 (Usage)

ブラウザで `http://localhost:5173` にアクセスします。

- **初回起動時**: 各社のセッション情報（`*_cookie.json`）が読み込まれ、自動的にポイントが同期されます。
- **キャッシュ**: 同期されたデータはブラウザの `localStorage` に保存され、次回アクセス時は瞬時に表示されます。
- **手動同期**: 右上の「Sync」ボタンを押すことで、最新の情報を取得します。

