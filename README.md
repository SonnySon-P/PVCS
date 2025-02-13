# VCS

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
