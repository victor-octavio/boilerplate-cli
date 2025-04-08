package cmd

import "fmt"

func generateFinalMessage() {
	cyan := "\033[36m"
	green := "\033[32m"

	// Arte ASCII e mensagem final em amarelo
	fmt.Println(string(green) + "    __          _ __                __      __                  ___ ")
	fmt.Println("   / /_  ____  (_) /__  _________  / /___ _/ /____        _____/ (_)")
	fmt.Println("  / __ \\/ __ \\/ / / _ \\/ ___/ __ \\/ / __ `/ __/ _ \\______/ ___/ / / ")
	fmt.Println(" / /_/ / /_/ / / /  __/ /  / /_/ / / /_/ / /_/  __/_____/ /__/ / /  ")
	fmt.Println("/_.___/\\____/_/_/\\___/_/  / .___/_/\\__,_/\\__/\\___/      \\___/_/_/   ")
	fmt.Println("                         /_/                                        ")
	fmt.Println()
	fmt.Println(string(cyan) + "\n Enjoy your new project 🚀 ")
	fmt.Println()
	fmt.Println()
}
