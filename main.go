package main

import (
	"VCSProject/vcs"
	"fmt"
	"os"
)

func main() {
	// 創建VCS
	vcs := vcs.NewVCS()

	// 檢查是否有action參數
	if len(os.Args) < 2 {
		fmt.Println("Error: action is required (init, add, remove, commit, status, log, checkout, create-branch, checkout-branch, merge)")
		return
	}

	// 根據action執行不同的邏輯
	switch os.Args[1] {
	case "init":
		err := vcs.Init()
		if err != nil {
			fmt.Println("Error:", err)
		}
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: add <filename>")
			return
		}
		err := vcs.Add(os.Args[2])
		if err != nil {
			fmt.Println("Error:", err)
		}
	case "remove":
		if len(os.Args) < 3 {
			fmt.Println("Usage: remove <filename>")
			return
		}
		err := vcs.Remove(os.Args[2])
		if err != nil {
			fmt.Println("Error:", err)
		}
	case "commit":
		if len(os.Args) < 3 {
			fmt.Println("Usage: commit <message>")
			return
		}
		err := vcs.Commit(os.Args[2])
		if err != nil {
			fmt.Println("Error:", err)
		}
	case "log":
		err := vcs.Log()
		if err != nil {
			fmt.Println("Error:", err)
		}
	case "status":
		err := vcs.Status()
		if err != nil {
			fmt.Println("Error:", err)
		}
	case "checkout":
		if len(os.Args) < 3 {
			fmt.Println("Usage: checkout <version number>")
			return
		}
		var version int
		_, err1 := fmt.Sscanf(os.Args[2], "%d", &version)
		if err1 != nil {
			fmt.Println("The version number format is incorrect.")
			return
		}
		err2 := vcs.Checkout(version)
		if err2 != nil {
			fmt.Println("Error:", err2)
		}
	case "create-branch":
		if len(os.Args) < 3 {
			fmt.Println("Usage: create-branch <branch name>")
			return
		}
		err := vcs.CreateBranch(os.Args[2])
		if err != nil {
			fmt.Println("Error:", err)
		}
	case "checkout-branch":
		if len(os.Args) < 3 {
			fmt.Println("Usage: checkout-branch <branch name>")
			return
		}
		err := vcs.CheckoutBranch(os.Args[2])
		if err != nil {
			fmt.Println("Error:", err)
		}
	case "merge":
		if len(os.Args) < 4 {
			fmt.Println("Usage: merge <target branch name> <source branch name>")
			return
		}
		targetBranch := os.Args[2]
		sourceBranch := os.Args[3]

		err := vcs.Merge(targetBranch, sourceBranch)
		if err != nil {
			fmt.Println("Error:", err)
		}
	default:
		fmt.Println("Error: invalid action. Choices are (init, add, remove, commit, status, log, checkout, create-branch, checkout-branch, merge)")
		return
	}
}
