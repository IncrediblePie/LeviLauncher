# 系統要求與安裝

本頁說明在正常安裝及管理 Minecraft Bedrock (GDK / UWP) 版本前，LeviLauncher 需要滿足哪些條件。

UWP 版本需要啟用 Windows 開發人員模式，並安裝套件聲明的 UWP 框架依賴項（例如 Microsoft.VCLibs）。下文 Gaming Services 與 GameInput 的要求適用於 GDK。GDK 版本以及 1.19.70.2 或以上的 UWP x64 套件支援實例隔離；未啟用隔離的 UWP 實例會共用該 Windows 應用程式的 `LocalState` 內容。UWP x64 實例亦支援主控台及相容的 `preload-native` DLL 模組，但 LeviLamina 自動安裝僅適用於受支援的 GDK 版本。

## 系統要求

| 項目 | 要求 |
| --- | --- |
| 作業系統 | Windows 10 或 Windows 11 |
| 遊戲版本 | Minecraft Bedrock Edition (GDK / UWP) |
| 授權 | 綁定在 Microsoft 帳號下的正版授權 |
| 網絡 | 用於下載版本、獲取元數據、測速鏡像與檢查更新 |

## 必要的 Windows 組件

首次啟動或安裝前，LeviLauncher 可能會引導你安裝缺失的組件。

- **Microsoft Gaming Services**（GDK）
- **Microsoft GameInput**（GDK）
- **WebView2 Runtime**
- **Microsoft Visual C++ 可轉散發套件（x64）**：原生模組載入器需要，LeviLauncher 會檢查並提示安裝

具體是否缺失，取決於你的 Windows 環境狀態。

### 使用自訂 WebView2 Fixed Version Runtime

如果系統中沒有 Evergreen WebView2 Runtime，但你已自行管理 WebView2 Fixed Version Runtime，可以使用以下任一方式啟動 LeviLauncher。

透過命令列參數指定：

```powershell
LeviLauncher.exe --webview2-runtime-dir "D:\WebView2\FixedRuntime"
```

或者在 `LeviLauncher.exe` 所在目錄建立 `webview2.json`：

```json
{
  "browserExecutableFolder": "D:\\WebView2\\FixedRuntime"
}
```

相對路徑以 `webview2.json` 所在目錄為基準。指定目錄必須直接包含 `msedgewebview2.exe`。命令列參數的優先次序高於設定檔；如果兩者都沒有設定，LeviLauncher 會繼續使用系統安裝的 WebView2 Runtime。

## 在安裝版本之前

請先完成這份檢查清單：

1. 至少從 Microsoft Store 安裝過一次 Minecraft Bedrock。
2. 如果商店狀態異常，先啟動一次遊戲以確認安裝完整。
3. 在使用 LeviLauncher 安裝或管理版本前，請先關閉遊戲。

## 安裝 LeviLauncher 主體

### 方案 A：GitHub Releases 頁面

適合希望直接從 LeviLauncher 官方下載頁獲取安裝程式，並順便查看更新記錄的用戶。

1. 開啟 LeviLauncher 的 [GitHub Releases](https://github.com/LiteLDev/LeviLauncher/releases) 頁面。
2. 下載安裝程式。
3. 執行並完成安裝嚮導。

### 方案 B：藍奏雲鏡像

如果你所在地區訪問 GitHub 速度較慢，這個入口通常更方便。

1. 開啟 [藍奏雲](https://levimc.lanzoue.com/b016ke39hc)。
2. 輸入密碼 `levi`。
3. 下載後在本地執行安裝程式。

## 安裝第一個託管版本

1. 在 LeviLauncher 內開啟 **下載**。
2. 選擇你要安裝的 Minecraft **正式版** 或 **預覽版**，以及 **GDK** 或 **UWP** 套件類型。
3. 選擇目標版本條目。
4. 決定是否啟用隔離。
5. 開始安裝並等待完成。

## 建議安裝策略

### 什麼時候選擇正式版（Release）

- 你想要更穩定的日常遊玩環境
- 你在經營長期世界
- 你希望模組與資源包的變動較少

### 什麼時候選擇預覽版（Preview）

- 你想提前體驗未來功能
- 你能接受不穩定性或相容性變化
- 你願意將 Minecraft 預覽版（Preview）環境與日常遊玩環境分開

::: tip 多數玩家的推薦做法
先建立一個**隔離的正式版（Release）版本**。只有在明確需要體驗預覽內容時，再額外添加**預覽版（Preview）**。
:::

## 如果安裝無法繼續

以下問題均可參考 [更新與故障排查](./update-troubleshooting)：

- 缺少 Gaming Services
- 缺少 GameInput
- 商店授權或安裝狀態不完整
- 目標路徑不可寫入
- 下載或鏡像失敗