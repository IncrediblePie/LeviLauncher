# システム要件とインストール

このページでは、LeviLauncher が Minecraft Bedrock (GDK / UWP) のバージョンを正常にインストールおよび管理する前に満たす必要がある条件について説明します。

UWP バージョンには、Windows の開発者モードと、パッケージが宣言する UWP フレームワーク依存関係（Microsoft.VCLibs など）が必要です。以下の Gaming Services と GameInput の要件は GDK に適用されます。インスタンスの隔離は GDK と、バージョン 1.19.70.2 以降の UWP x64 パッケージで利用できます。隔離されていない UWP インスタンスは、Windows アプリの `LocalState` コンテンツを共有します。UWP x64 インスタンスはコンソールと互換性のある `preload-native` DLL Mod にも対応していますが、LeviLamina の自動インストールはサポートされている GDK バージョンでのみ利用できます。

## システム要件

| 項目 | 要件 |
| --- | --- |
| OS | Windows 10 または Windows 11 |
| ゲームバージョン | Minecraft Bedrock Edition (GDK / UWP) |
| ライセンス | Microsoft アカウントに紐づけられた正規版ライセンス |
| ネットワーク | バージョンのダウンロード、メタデータの取得、ミラーの速度測定、更新の確認に使用 |

## 必須の Windows コンポーネント

初回起動時またはインストール前に、LeviLauncher は不足しているコンポーネントのインストールを促す場合があります。

- **Microsoft Gaming Services** (GDK)
- **Microsoft GameInput** (GDK)
- **WebView2 Runtime**
- **Microsoft Visual C++ Redistributable (x64)**：ネイティブ Mod ローダーに必要です。LeviLauncher が有無を確認し、インストールを提案します

具体的にどのコンポーネントが不足しているかは、あなたの Windows 環境の状態によって異なります。

### カスタムの WebView2 Fixed Version Runtime を使用する

Evergreen WebView2 Runtime がインストールされておらず、WebView2 Fixed Version Runtime を自分で管理している場合は、次のいずれかの方法で LeviLauncher を起動してください。

コマンドラインでディレクトリを指定します：

```powershell
LeviLauncher.exe --webview2-runtime-dir "D:\WebView2\FixedRuntime"
```

または、`LeviLauncher.exe` と同じ場所に `webview2.json` を作成します：

```json
{
  "browserExecutableFolder": "D:\\WebView2\\FixedRuntime"
}
```

相対パスは `webview2.json` があるディレクトリを基準に解決されます。選択したディレクトリには `msedgewebview2.exe` が直接含まれている必要があります。コマンドラインオプションは構成ファイルより優先されます。どちらも設定されていない場合、LeviLauncher は引き続きシステムにインストールされている WebView2 Runtime を使用します。

## バージョンをインストールする前に

以下のチェックリストを完了してください：

1. Microsoft Store から少なくとも一度 Minecraft Bedrock をインストール済みであることを確認します。
2. Store の状態が不安定な場合は、一度ゲームを起動してインストールが完了していることを確認します。
3. LeviLauncher を使用してバージョンをインストールまたは管理する前に、ゲームを終了しておいてください。

## LeviLauncher 本体のインストール

### 方針 A：GitHub Releases ページ

公式ダウンロードページから直接インストーラーを取得し、更新履歴も同時に確認したいユーザー向けです。

1. LeviLauncher の [GitHub Releases](https://github.com/LiteLDev/LeviLauncher/releases) ページを開きます。
2. インストーラーをダウンロードします。
3. インストールウィザードを実行して完了します。

### 方針 B：藍奏雲ミラー

GitHub へのアクセス速度が遅い地域のユーザーには、こちらのエントリーポイントが便利です。

1. [藍奏雲](https://levimc.lanzoue.com/b016ke39hc) を開きます。
2. パスワード `levi` を入力します。
3. ローカルでインストーラーを実行します。

## 最初のホストバージョンのインストール

1. LeviLauncher 内で **ダウンロード** を開きます。
2. インストールする Minecraft のビルドとして **リリース** または **プレビュー** を、パッケージの種類として **GDK** または **UWP** を選択します。
3. 対象のバージョンエントリを選択します。
4. 隔離を有効にするかどうかを決定します。
5. インストールを開始し、完了まで待ちます。

## 推奨されるインストール戦略

### いつ正式版（Release）を選ぶべきか

- より安定した日常プレイ環境が欲しい場合
- 長期的なワールドを運営している場合
- Mod やリソースパックの変更を少なくしたい場合

### いつプレビュー版（Preview）を選ぶべきか

- 未来の機能を先行体験したい場合
- 不安定さや互換性の変更を許容できる場合
- Minecraft プレビュー版（Preview）環境と日常プレイ環境を分離したい場合

::: tip 多くのプレイヤーに推奨される方法
まず **隔離された正式版（Release）バージョン** を構築します。プレビュー内容を明確に体験したい場合に限り、**プレビュー版（Preview）** を追加でインストールすることを推奨します。
:::

## インストールが続行できない場合

以下の問題は、[更新とトラブルシューティング](./update-troubleshooting) を参照してください：

- Gaming Services が不足している
- GameInput が不足している
- Store のライセンスまたはインストール状態が不完全
- 対象パスに書き込み権限がない
- ダウンロードまたはミラーに失敗