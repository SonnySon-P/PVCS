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
