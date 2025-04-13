package vcs

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// VCS資料結構
type VCS struct {
	repoDirectory    string
	filesDirectory   string
	historyDirectory string
	currentBranch    string
	currentVersion   int
}

// 創建VCS
func NewVCS() *VCS {
	repoDirectory := ".vcs" // 隱藏的資料夾，用來儲存各版本檔案
	filesDirectory := filepath.Join(repoDirectory, "files")
	historyDirectory := filepath.Join(repoDirectory, "history")
	currentBranch := "main"
	currentVersion := 0
	return &VCS{repoDirectory: repoDirectory, filesDirectory: filesDirectory, historyDirectory: historyDirectory, currentBranch: currentBranch, currentVersion: currentVersion}
}

// 初始化VCS，創建必要的文件夹
func (v *VCS) Init() error {
	// 檢查repoDirectory是否存在
	if _, err1 := os.Stat(v.repoDirectory); os.IsNotExist(err1) {
		// 創建filesDirectory
		err2 := os.MkdirAll(v.filesDirectory, os.ModePerm)
		if err2 != nil {
			return fmt.Errorf("unable to create document folder: %v", err2)
		}

		// 創建historyDirectory
		err3 := os.MkdirAll(v.historyDirectory, os.ModePerm)
		if err3 != nil {
			return fmt.Errorf("unable to create history folder: %v", err3)
		}

		// 創建main branch資料夹
		mainBranchDirectory := filepath.Join(v.historyDirectory, "main")
		err4 := os.MkdirAll(mainBranchDirectory, os.ModePerm)
		if err4 != nil {
			return fmt.Errorf("unable to create main branch folder: %v", err4)
		}

		// 將currentBranch寫入檔案
		err5 := v.writeCurrentBranch()
		if err5 != nil {
			return nil
		}

		fmt.Printf("Initialized empty VCS repository in %s\n", v.repoDirectory)
		return nil
	} else {
		return fmt.Errorf("repository already exists at %s", v.repoDirectory)
	}
}

// 將文件添加到版本控制
func (v *VCS) Add(filename string) error {
	// 檢查files資料夾是否存在
	if _, err1 := os.Stat(v.filesDirectory); os.IsNotExist(err1) {
		err2 := os.MkdirAll(v.filesDirectory, os.ModePerm)
		if err2 != nil {
			return fmt.Errorf("unable to create files directory: %v", err2)
		}
	}

	// 檢查檔案是否存在
	if _, err3 := os.Stat(filename); os.IsNotExist(err3) {
		return fmt.Errorf("file does not exist: %s", filename)
	}

	// 複製檔案到file資料夾
	destinationPath := filepath.Join(v.filesDirectory, filepath.Base(filename))
	err4 := copyFile(filename, destinationPath)
	if err4 != nil {
		return fmt.Errorf("failed to add file: %v", err4)
	}

	fmt.Printf("Added %s to version control.\n", filename)
	return nil
}

// 將files中的指定資料夾或檔案移除
func (v *VCS) Remove(filename string) error {
	removePath := filepath.Join(v.filesDirectory, filepath.Base(filename))

	// 檢查路徑是否存在
	info, err1 := os.Stat(removePath)
	if err1 != nil {
		return fmt.Errorf("cannot access path %s: %v", removePath, err1)
	}

	// 檢查路徑是否為資料夾
	if info.IsDir() {
		// 移除資料夾
		err2 := os.RemoveAll(removePath)
		if err2 != nil {
			return fmt.Errorf("cannot delete directory %s: %v", removePath, err2)
		}
		fmt.Printf("Directory %s has been successfully deleted.\n", removePath)
	} else {
		// 移除檔案
		err3 := os.Remove(removePath)
		if err3 != nil {
			return fmt.Errorf("cannot delete file %s: %v", removePath, err3)
		}
		fmt.Printf("File %s has been successfully deleted.\n", removePath)
	}
	return nil
}

// 提交目前狀態，並產生新版本
func (v *VCS) Commit(message string) error {
	// 從檔案讀取currentBranch
	err1 := v.readCurrentBranch()
	if err1 != nil {
		return err1
	}

	// 將版本編號加1
	v.currentVersion = v.getCurrentVersionOfBranch(v.currentBranch) + 1

	// 生成一個新版本的路徑
	versionDirectory := filepath.Join(v.historyDirectory, v.currentBranch, fmt.Sprintf("version_%d", v.currentVersion))

	// 創建新版本資料夾
	err2 := os.Mkdir(versionDirectory, os.ModePerm)
	if err2 != nil {
		return fmt.Errorf("unable to create version folder: %v", err2)
	}

	// 複製檔案到版本資料夾
	files, err3 := os.ReadDir(v.filesDirectory)
	if err3 != nil {
		return fmt.Errorf("unable to read folder: %v", err3)
	}

	for _, file := range files {
		filePath := filepath.Join(v.filesDirectory, file.Name())
		if info, err4 := os.Stat(filePath); err4 == nil && !info.IsDir() {
			err5 := copyFile(filePath, filepath.Join(versionDirectory, file.Name()))
			if err5 != nil {
				return fmt.Errorf("file copy failure: %v", err5)
			}
		}
	}

	// 更新目前version為新version
	err6 := v.writeCurrentVersion()
	if err6 != nil {
		return err6
	}

	// 寫入提交訊息
	messagePath := filepath.Join(versionDirectory, "commit_message.txt")
	err7 := os.WriteFile(messagePath, []byte(message), 0644)
	if err7 != nil {
		return fmt.Errorf("failed to write commit message: %v", err7)
	}

	fmt.Printf("Committed version %d with message: %s\n", v.currentVersion, message)
	return nil
}

// 取得branch所有提交記錄
func (v *VCS) Log() error {
	// 從檔案讀取currentBranch
	err1 := v.readCurrentBranch()
	if err1 != nil {
		return err1
	}

	fmt.Printf("On the %s branch\n", v.currentBranch)

	// 只查看目前branch的所有版本
	branchDirectory := filepath.Join(v.historyDirectory, v.currentBranch)
	files, err2 := os.ReadDir(branchDirectory)
	if err2 != nil {
		return fmt.Errorf("unable to read folder: %v", err2)
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "version_") {
			version := strings.TrimPrefix(file.Name(), "version_")
			messagePath := filepath.Join(branchDirectory, file.Name(), "commit_message.txt")
			message, err3 := os.ReadFile(messagePath)
			if err3 != nil {
				return fmt.Errorf("unable to read commit message for version %s: %v", version, err3)
			}
			fmt.Printf("Version %s: %s\n", version, message)
		}
	}
	return nil
}

// 狀態查看，搜尋結果目前資料夾內容
func (v *VCS) Status() error {
	// 從檔案讀取currentBranch
	err1 := v.readCurrentBranch()
	if err1 != nil {
		return err1
	}

	// 從檔案讀取currentVersion
	err2 := v.readCurrentVersion()
	if err2 != nil {
		return err2
	}

	fmt.Printf("On the %s branch, version %d\n", v.currentBranch, v.currentVersion)

	files, err := os.ReadDir(v.filesDirectory)
	if err != nil {
		return fmt.Errorf("unable to read folder: %v", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no files are being tracked")
	}

	// 若有檔案，列出追蹤的檔案
	fmt.Println("Tracked files:")
	for _, file := range files {
		fmt.Println(file.Name())
	}
	return nil
}

// 切換到指定版本
func (v *VCS) Checkout(version int) error {
	// 從檔案讀取currentBranch
	err1 := v.readCurrentBranch()
	if err1 != nil {
		return nil
	}

	versionDirectory := filepath.Join(v.historyDirectory, v.currentBranch, fmt.Sprintf("version_%d", version))
	files, err2 := os.ReadDir(versionDirectory)
	if err2 != nil {
		return fmt.Errorf("version %d does not exist: %v", version, err2)
	}

	// 取得.vcs的父資料夾，開發程式所在的工作目錄
	workingDirectory := filepath.Dir(v.repoDirectory)

	// 複製檔案
	for _, file := range files {
		if file.Name() != "commit_message.txt" {
			sourcePath := filepath.Join(versionDirectory, file.Name())

			// 清空暫存區資料夾
			err3 := clearFolder(v.filesDirectory)
			if err3 != nil {
				return err3
			}

			// 複製到工作區
			destinationPathRepo := filepath.Join(workingDirectory, file.Name())
			err4 := copyFile(sourcePath, destinationPathRepo)
			if err4 != nil {
				return fmt.Errorf("unable to switch workspace version for %s: %v", file.Name(), err4)
			}

			// 複製到暫存區
			destinationPath := filepath.Join(v.filesDirectory, file.Name())
			err5 := copyFile(sourcePath, destinationPath)
			if err5 != nil {
				return fmt.Errorf("unable to switch the staging area version for %s: %v", file.Name(), err5)
			}
		}
	}

	// 更新目前version為新version
	v.currentVersion = version
	err5 := v.writeCurrentVersion()
	if err5 != nil {
		return err5
	}

	fmt.Printf("Checked out version %d\n", version)
	return nil
}

// 創建分支
func (v *VCS) CreateBranch(branchName string) error {
	// 從檔案讀取currentBranch
	err1 := v.readCurrentBranch()
	if err1 != nil {
		return err1
	}

	branchDirectory := filepath.Join(v.historyDirectory, branchName)
	if _, err2 := os.Stat(branchDirectory); !os.IsNotExist(err2) {
		return fmt.Errorf("branch %s already exists", branchName)
	}

	// 創建branch資料夾
	err3 := os.Mkdir(branchDirectory, os.ModePerm)
	if err3 != nil {
		return fmt.Errorf("unable to create branch directory: %v", err3)
	}

	// 將目前版本的檔案複製到新branch
	v.currentVersion = v.getCurrentVersionOfBranch(v.currentBranch)
	sourceVersionDirectory := filepath.Join(v.historyDirectory, v.currentBranch, fmt.Sprintf("version_%d", v.currentVersion))
	destinationVersionDirectory := filepath.Join(branchDirectory, fmt.Sprintf("version_%d", v.currentVersion))
	err4 := os.Mkdir(destinationVersionDirectory, os.ModePerm)
	if err4 != nil {
		return fmt.Errorf("unable to create version folder: %v", err4)
	}
	err5 := v.createVersionSnapshot(sourceVersionDirectory, destinationVersionDirectory)
	if err5 != nil {
		return fmt.Errorf("unable to create snapshot for branch: %v", err5)
	}

	// 更新目前branch為新branch
	v.currentBranch = branchName
	err6 := v.writeCurrentBranch()
	if err6 != nil {
		return err6
	}

	// 更新目前version為新version
	err7 := v.writeCurrentVersion()
	if err7 != nil {
		return err7
	}

	fmt.Printf("Branch %s created successfully\n", branchName)
	return nil
}

// 切換branch
func (v *VCS) CheckoutBranch(branchName string) error {
	branchDirectory := filepath.Join(v.historyDirectory, branchName)
	if _, err1 := os.Stat(branchDirectory); os.IsNotExist(err1) {
		return fmt.Errorf("branch %s does not exist", branchName)
	}

	// 更新目前分支為指定branch
	v.currentBranch = branchName
	err2 := v.writeCurrentBranch()
	if err2 != nil {
		return err2
	}

	v.currentVersion = v.getCurrentVersionOfBranch(v.currentBranch)
	v.Checkout(v.currentVersion)

	// 更新目前version為新version
	err3 := v.writeCurrentVersion()
	if err3 != nil {
		return err3
	}

	fmt.Printf("Checked out to branch: %s\n", branchName)
	return nil
}

// 合併來源branch到目標branch
func (v *VCS) Merge(targetBranch, sourceBranch string) error {
	// 目標branch或來源branch不存在，則傳回錯誤
	targetBranchDirectory := filepath.Join(v.historyDirectory, targetBranch)
	sourceBranchDirectory := filepath.Join(v.historyDirectory, sourceBranch)

	if _, err1 := os.Stat(targetBranchDirectory); err1 != nil {
		if os.IsNotExist(err1) {
			return fmt.Errorf("target branch %s does not exist", targetBranch)
		}
		return fmt.Errorf("error checking target branch %s: %v", targetBranch, err1)
	}

	if _, err2 := os.Stat(sourceBranchDirectory); err2 != nil {
		if os.IsNotExist(err2) {
			return fmt.Errorf("source branch %s does not exist", sourceBranch)
		}
		return fmt.Errorf("error checking source branch %s: %v", sourceBranch, err2)
	}

	// 取得目標branch與來源branch的最大版本路徑
	targetVersion := v.getCurrentVersionOfBranch(targetBranch)
	sourceVersion := v.getCurrentVersionOfBranch(sourceBranch)

	targetVersionDirectory := filepath.Join(targetBranchDirectory, fmt.Sprintf("version_%d", targetVersion))
	targetFiles, err3 := os.ReadDir(targetVersionDirectory)
	if err3 != nil {
		return fmt.Errorf("failed to read target branch version file: %s", err3)
	}

	sourceVersionDirectory := filepath.Join(sourceBranchDirectory, fmt.Sprintf("version_%d", sourceVersion))
	sourceFiles, err4 := os.ReadDir(sourceVersionDirectory)
	if err4 != nil {
		return fmt.Errorf("failed to read source branch version file: %s", err3)
	}

	// 將來源branch檔案合併到目標branch
	mergeVersionDirectory := filepath.Join(targetBranchDirectory, fmt.Sprintf("version_%d", targetVersion+1))
	err5 := os.Mkdir(mergeVersionDirectory, os.ModePerm)
	if err5 != nil {
		return fmt.Errorf("unable to create merged revision folder: %s", err4)
	}

	// 複製文件
	for _, targetFile := range targetFiles {
		if targetFile.Name() != "commit_message.txt" {
			targetFilePath := filepath.Join(targetVersionDirectory, targetFile.Name())

			// 問使用者是否要複製檔案
			fmt.Printf("Target File: %s\n", targetFilePath)
			var userChoice string
			fmt.Print("Do you want to copy this file to the merge directory? (yes/no): ")
			fmt.Scanln(&userChoice)

			// 檢查使用者輸入
			if strings.ToLower(userChoice) == "yes" {
				destinationPath := filepath.Join(mergeVersionDirectory, targetFile.Name())
				err5 := copyFile(targetFilePath, destinationPath)
				if err5 != nil {
					return fmt.Errorf("failed to copy file: %s", err5)
				}
			}
		}
	}

	for _, sourceFile := range sourceFiles {
		if sourceFile.Name() != "commit_message.txt" {
			sourceFilePath := filepath.Join(targetVersionDirectory, sourceFile.Name())

			// 問使用者是否要複製檔案
			fmt.Printf("Source File: %s\n", sourceFilePath)
			var userChoice1 string
			fmt.Print("Do you want to copy this file to the merge directory? (yes/no): ")
			fmt.Scanln(&userChoice1)

			// 檢查使用者輸入
			if strings.ToLower(userChoice1) == "yes" {
				_, err6 := os.Stat(sourceFilePath)
				if os.IsNotExist(err6) {
					destinationPath := filepath.Join(mergeVersionDirectory, sourceFile.Name())
					err7 := copyFile(sourceFilePath, destinationPath)
					if err7 != nil {
						return fmt.Errorf("failed to copy file: %s", err7)
					}
				} else {
					var userChoice2 string
					fmt.Print("Do you want to overwrite the target file?? (yes/no): ")
					fmt.Scanln(&userChoice2)
					if strings.ToLower(userChoice2) == "yes" {
						_, err8 := os.Stat(sourceFilePath)
						if os.IsNotExist(err8) {
							destinationPath := filepath.Join(mergeVersionDirectory, sourceFile.Name())
							err9 := copyFile(sourceFilePath, destinationPath)
							if err9 != nil {
								return fmt.Errorf("failed to copy file: %s", err9)
							}
						}
					}
				}
			}
		}
	}

	// 合併完成，提交訊息
	commitMessage := fmt.Sprintf("Merged %s into %s", sourceBranch, targetBranch)
	messagePath := filepath.Join(mergeVersionDirectory, "commit_message.txt")
	err6 := os.WriteFile(messagePath, []byte(commitMessage), 0644)
	if err6 != nil {
		return fmt.Errorf("failed to write commit message: %s", err6)
	}

	// 更新目前分支為指定branch
	v.currentBranch = targetBranch
	err7 := v.writeCurrentBranch()
	if err7 != nil {
		return err7
	}

	// 更新目前version為新version
	v.currentVersion = targetVersion + 1
	err8 := v.writeCurrentVersion()
	if err8 != nil {
		return err8
	}

	// 提交合併後的版本
	fmt.Printf("Successfully merged %s into %s\n", sourceBranch, targetBranch)
	return nil
}

// 紀錄目前branch
func (v *VCS) writeCurrentBranch() error {
	// 在.vcs資料夾中創建並寫入currentBranch
	currentBranchFilePath := filepath.Join(v.repoDirectory, "currentBranch.txt")
	err := os.WriteFile(currentBranchFilePath, []byte(v.currentBranch), 0644)
	if err != nil {
		return fmt.Errorf("unable to write current branch file: %v", err)
	}
	return nil
}

// 載入目前branch
func (v *VCS) readCurrentBranch() error {
	currentBranchFilePath := filepath.Join(v.repoDirectory, "currentBranch.txt")

	// 檢查檔案是否存在
	if _, err := os.Stat(currentBranchFilePath); os.IsNotExist(err) {
		return fmt.Errorf("current branch file does not exist")
	}

	// 讀取當前branch名稱
	branchName, err := os.ReadFile(currentBranchFilePath)
	if err != nil {
		return fmt.Errorf("unable to read current branch file: %v", err)
	}
	v.currentBranch = string(branchName)
	return nil
}

// 紀錄目前版本
func (v *VCS) writeCurrentVersion() error {
	// 在.vcs資料夾中創建並寫入currentVersion
	currentVersionFilePath := filepath.Join(v.repoDirectory, "currentVersion.txt")
	err := os.WriteFile(currentVersionFilePath, []byte(strconv.Itoa(v.currentVersion)), 0644)
	if err != nil {
		return fmt.Errorf("unable to write current branch file: %v", err)
	}
	return nil
}

// 載入目前版本
func (v *VCS) readCurrentVersion() error {
	currentVersionFilePath := filepath.Join(v.repoDirectory, "currentVersion.txt")

	// 檢查檔案是否存在
	if _, err1 := os.Stat(currentVersionFilePath); os.IsNotExist(err1) {
		return fmt.Errorf("current branch file does not exist")
	}

	// 讀取當前版本編號
	versionNumber, err2 := os.ReadFile(currentVersionFilePath)
	if err2 != nil {
		return fmt.Errorf("unable to read current branch file: %v", err2)
	}
	currentVersion, err3 := strconv.Atoi(string(versionNumber))
	if err3 != nil {
		return fmt.Errorf("conversion error: %v", err3)
	}
	v.currentVersion = currentVersion
	return nil
}

// 取得指定分支的目前版本
func (v *VCS) getCurrentVersionOfBranch(branch string) int {
	branchDirectory := filepath.Join(v.historyDirectory, branch)
	files, _ := os.ReadDir(branchDirectory)
	versionNumbers := []int{}
	for _, file := range files {
		// 檢查檔名是否有version_的前綴開頭
		if strings.HasPrefix(file.Name(), "version_") {
			var version int
			// 取得檔案名稱中的版本
			fmt.Sscanf(file.Name(), "version_%d", &version)
			versionNumbers = append(versionNumbers, version)
		}
	}

	// 找不到就回傳0
	if len(versionNumbers) == 0 {
		return 0
	}

	// 回傳最大的版本
	maxVersion := versionNumbers[0]
	for _, v := range versionNumbers {
		if v > maxVersion {
			maxVersion = v
		}
	}
	return maxVersion
}

// 清空資料夾
func clearFolder(folderPath string) error {
	// 取得資料夾中的所有檔案
	files, err1 := os.ReadDir(folderPath)
	if err1 != nil {
		return fmt.Errorf("unable to read folder: %v", err1)
	}

	// 遍歷資料夾中的每個檔案
	for _, file := range files {
		filePath := filepath.Join(folderPath, file.Name())
		// 刪除檔案
		err2 := os.Remove(filePath)
		if err2 != nil {
			return fmt.Errorf("failed to delete file %s: %v", filePath, err2)
		}
	}
	return nil
}

// 複製文件
func copyFile(sourcePath, destinationPath string) error {
	// 讀取源文件
	input, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("unable to read source file %s: %v", sourcePath, err)
	}

	// 寫入到目標文件
	err = os.WriteFile(destinationPath, input, 0644)
	if err != nil {
		return fmt.Errorf("unable to write to destination file %s: %v", destinationPath, err)
	}
	return nil
}

// 建立branch的版本快照
func (v *VCS) createVersionSnapshot(sourceVersionDirectory, destinationVersionDirectory string) error {
	files, err1 := os.ReadDir(sourceVersionDirectory)
	if err1 != nil {
		return fmt.Errorf("unable to read version directory: %v", err1)
	}

	// 清空暫存區資料夾
	err2 := clearFolder(v.filesDirectory)
	if err2 != nil {
		return err2
	}

	// 取得.vcs的父資料夾，開發程式所在的工作目錄
	workingDirectory := filepath.Dir(v.repoDirectory)

	// 建立branch目錄下的版本快照
	for _, file := range files {
		// 將檔案複製到新的branch
		sourcePath := filepath.Join(sourceVersionDirectory, file.Name())
		destinationPath := filepath.Join(destinationVersionDirectory, file.Name())
		err3 := copyFile(sourcePath, destinationPath)
		if err3 != nil {
			return fmt.Errorf("unable to copy file to branch: %v", err3)
		}

		if file.Name() != "commit_message.txt" {
			// 複製到工作區
			destinationPathRepo := filepath.Join(workingDirectory, file.Name())
			err4 := copyFile(sourcePath, destinationPathRepo)
			if err4 != nil {
				return fmt.Errorf("unable to switch workspace version for %s: %v", file.Name(), err4)
			}

			// 複製到暫存區
			destinationPathFile := filepath.Join(v.filesDirectory, file.Name())
			err5 := copyFile(sourcePath, destinationPathFile)
			if err5 != nil {
				return fmt.Errorf("unable to switch the staging area version for %s: %v", file.Name(), err5)
			}
		}
	}
	return nil
}
