# Simple Version Control Software

基於Linux環境開發的簡單防毒毒軟體框架(無病毒碼)。

## 壹、基本說明
**一、目標：**
回想起剛開始接觸Windows作業系統的那段時光，每當我們完成作業系統的安裝後，似乎總少不了的一個步驟便是安裝防毒軟體。那時，我對防毒軟體的運作原理充滿了好奇，尤其是它如何有效地保護我們的系統免受各種威脅。這份好奇心驅使我深入思考防毒軟體背後的運作邏輯。藉著這次的寫作機會，我希望能夠簡單地梳理並釐清防毒軟體的基本架構，進一步了解其如何運行、偵測與防範病毒的過程，並探索其在保障數位安全方面所扮演的關鍵角色。
<br>

**二、開發環境：**
以下是開發該平台所採用的環境：
* 虛擬機：Docker
* 映像檔：golang
* 程式語言：Golang
* 程式編輯器：Visual Studio Code

**三、檔案說明：** 
此專案檔案（指coding這個資料夾）主要分為兩個資料夾：nodejs和tests。其中，nodejs資料夾為後端平台的主要程式碼，tests資料夾則存放使用jest框架進行的單元測試。接下來將對各資料夾中的檔案內容進行詳細說明。
```bash
.
├── LICENSE
├── README.md
├── go.mod
├── main.go  # 主程式
└──  vcs
      └── vcs.go  # 各功能副程式
```

## 貳、操作說明
**一、安裝程式方式：** 
將一個編譯好的執行檔放置到`bin`資料夾並設置好環境路徑，步驟如下：
***步驟1: 編譯Golang，生成一個名為`vm`的執行檔。
```bash
go build -o vm main.go
```

***步驟 2: 將執行檔放到`bin`資料夾
請在UNIX類系統(如Linux或macOS)中，將執行檔放到`/usr/local/bin`或`~/bin` 通常是用來存放可執行檔的目錄。
```bash
mkdir -p ~/bin
mv vm ~/bin/
```

***步驟 3: 設置環境路徑
接下來，您需要設置您的環境變數，使得系統可以找到您放置的`bin`資料夾。使得可以在任何地方執行`vm`。

1. 開啟`.bashrc`或`.zshrc`配置檔(取決於您使用的 shell)，如果使用的是`bash`，需要編輯`~/.bashrc`文件；如果使用的是`zsh`，則是`~/.zshrc`文件。
```bash
nano ~/.bashrc  # 如果是 bash
# 或者
nano ~/.zshrc   # 如果是 zsh
```

2. 添加`bin`資料夾到`$PATH`，在配置文件中，加入以下一行：
```bash
export PATH=$PATH:~/bin
```

**二、運行程式方式：**
```bash
vm -load <path>  #  載入映像檔，將指定的映像檔解壓到容器目錄中
vm -save <path>  #  將當前容器保存為新映像檔
vm -create <path>  #  從映像檔創建容器
vm -start <containerID>  #  啟動容器
vm -stop <containerID>  #  停止容器
vm -delete <containerID>  #  刪除容器
```


# VCS

專案結構應該像這樣：
your_project_folder/
  ├── main.go
  └── vcs/
      └── vcs.go

repoDir(工作區)
  ├── filesDir(暫存區)
  └── historyDir(提交區)
        ├── main
        └── branch
              ├── version_1
              └── version_2
              
go run main.go init        # 初始化版本控制
go run main.go commit "Initial commit"  # 提交更改
go run main.go status      # 查看当前文件状态
go run main.go log         # 查看提交记录
go run main.go checkout 1  # 切换到版本 1

你提到的問題確實是一個很常見的挑戰。逐行比較的方式對於文件內容的微小變化（如空白行、換行等）非常敏感，這會導致無法識別相同的內容。尤其是在檔案中插入或刪除空白行時，可能會導致程式碼不一致的錯誤。

為了解決這個問題，我們可以透過以下幾種方法來改進：

1. 忽略空白行和空白字符
我們可以忽略空行和多餘的空格，這樣即使程式碼內容在結構上有所變化，只要內容本身沒變化，依然可以視為相同。

2. 使用行內差異計算
如果我們偵測到某一行相同，但該行的內容在空格或格式上有所不同，我們可以嘗試偵測 行內差異，並且根據該差異來更新檔案內容。

3. 利用最小編輯距離演算法（Levenshtein Distance）
另一個解決方案是採用 最小編輯距離演算法（Levenshtein Distance），用於計算兩個字串的差異，並能夠提供更細微的比較。這適用於行內容相似，但具體字元不同的情況。

我將為你提供一個簡化版的解決方案，其中考慮了忽略空行以及去除每行兩端的空白字元。



在 VCS 结构中增加分支管理的功能，如创建分支、切换分支、合并分支等。
分支管理将建立在 historyDir 目录内为每个分支创建独立的文件夹，每个分支的文件夹中会包含该分支的所有版本。
分支切换通过更新 filesDir 的路径，切换到当前分支的文件夹。
分支合并则是将源分支中的文件合并到目标分支中。

repoDir（工作區）：
這是開發者直接編輯和修改文件的地方。
開發者在此進行程式碼開發、文件編輯等操作。
repoDir 存放的是目前的工作文件。
filesDir（暫存區）：
當開發者想要將檔案加入版本控制時，檔案會先被放入 filesDir。
filesDir 是暫存區（Staging Area），存放的是準備提交的文件。
使用 add 指令時，檔案會從工作區 (repoDir) 複製到 filesDir 中。
historyDir（歷史資料夾）：
存放已經提交的版本文件。每次執行 commit 指令時，filesDir 中的檔案會被複製到 historyDir 中，並形成一個版本快照。
historyDir 保留了所有提交的歷史記錄，每個版本的文件都有一個獨立的目錄，包含版本內容和提交資訊。
工作流程：
新增檔案（add）：
開發者在 repoDir 中編輯檔案後，使用 add 指令將檔案新增至 filesDir。
add 並不會提交文件，它只是把文件從工作區（repoDir）移到暫存區（filesDir）。
提交文件（commit）：
當檔案已經在 filesDir 中時，執行 commit 操作會將 filesDir 中的檔案（包括提交資訊）複製到 historyDir 中，形成一個新的版本。
每次提交都會產生一個新的版本快照，保存在 historyDir 中。
切換版本（checkout）：
切換版本時，系統會將 historyDir 中某一版本的檔案複製回 repoDir（工作區）。
repoDir 會更新為目標版本的檔案狀態。
如何切换版本或分支时的操作流程
假设您需要切换到某个版本或分支时，您的 repoDir 会更新，以便包含目标版本（或分支）的文件内容。目标是确保工作区始终与版本控制系统中的文件一致。

分支合併（merge）是版本控制中非常重要的操作，通常用于将两个分支的更改合并到一个单一的分支中。在实现简单的 VCS 系统时，分支合并主要涉及以下几个步骤：

1. 理解分支结构
在当前的分支模型中，historyDir/{branch_name} 是存放每个分支历史版本的地方，而每个分支有独立的版本控制。合併操作需要将两个分支中的更改合并在一起。

2. 分支合併的基本逻辑
在最简单的 VCS 系统中，分支合并操作可以简单地理解为：将一个分支的内容合并到另一个分支。

选择一个目标分支（targetBranch），并从另一个分支（sourceBranch）中合并更改。
通过将 sourceBranch 的版本文件和提交信息复制到 targetBranch 来合并更改。
然后可以提交一个新的合并版本，将两者的内容合并到目标分支中。
3. 实现分支合併功能
下面是一个简单的分支合併示例。我们将根据目标分支（targetBranch）和源分支（sourceBranch）的版本内容来合并文件。
