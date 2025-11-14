package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx         context.Context
	title       string
	currentFile string
	content     string
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.title = "未命名文件"
	AppMenu := menu.NewMenu()
	item := AppMenu.AddSubmenu("文件")
	item.AddText("新建", keys.CmdOrCtrl("O"), func(cd *menu.CallbackData) {
		a.NewFile()
	})

	item.AddText("打开", keys.CmdOrCtrl("O"), func(cd *menu.CallbackData) {
		a.OpenFile()
	})

	item.AddText("保存", keys.CmdOrCtrl("S"), func(cd *menu.CallbackData) {
		a.SaveFile()
	})

	item.AddText("另存为", keys.CmdOrCtrl("E"), func(cd *menu.CallbackData) {
		a.SaveAsFile()
	})

	item.AddText("退出", keys.CmdOrCtrl("Q"), func(cd *menu.CallbackData) {
		runtime.Quit(a.ctx)
	})

	item = AppMenu.AddSubmenu("模式切换")
	item.AddText("普通模式", keys.CmdOrCtrl("1"), func(cd *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, "changeMode", "wysiwyg")

	})
	item.AddText("高级模式", keys.CmdOrCtrl("2"), func(cd *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, "changeMode", "markdown")
	})
	runtime.MenuSetApplicationMenu(a.ctx, AppMenu)
}

// OpenFile 打开文件并返回文件内容
func (a *App) OpenFile() (string, error) {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择Markdown文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Markdown文件 (*.md)",
				Pattern:     "*.md",
			},
			{
				DisplayName: "文本文件 (*.txt)",
				Pattern:     "*.txt",
			},
			{
				DisplayName: "所有文件 (*.*)",
				Pattern:     "*.*",
			},
		},
	})
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "打开文件错误",
			Message: err.Error(),
		})
		return "", err
	}

	if selection == "" {
		return "", fmt.Errorf("未选择文件")
	}

	// 读取文件内容
	content, err := os.ReadFile(selection)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "读取文件错误",
			Message: err.Error(),
		})
		return "", err
	}

	a.currentFile = selection
	a.content = string(content)
	a.title = strings.TrimSuffix(filepath.Base(selection), filepath.Ext(selection))
	runtime.EventsEmit(a.ctx, "openFile", a.content)

	return a.content, nil
}

// SaveFile 保存文件
func (a *App) SaveFile() (bool, error) {
	if a.currentFile == "" {
		return a.SaveAsFile()
	}

	err := os.WriteFile(a.currentFile, []byte(a.content), 0644)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "保存文件错误",
			Message: err.Error(),
		})
		return false, err
	}
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   "保存成功",
		Message: fmt.Sprintf("文件已保存到: %s", a.currentFile),
	})
	return true, nil
}

// SaveAsFile 另存为文件
func (a *App) SaveAsFile() (bool, error) {
	selection, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title: "保存文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Markdown文件 (*.md)",
				Pattern:     "*.md",
			},
			{
				DisplayName: "文本文件 (*.txt)",
				Pattern:     "*.txt",
			},
		},
		DefaultFilename: "untitled.md",
	})
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "保存文件错误",
			Message: err.Error(),
		})
		return false, err
	}

	if selection == "" {
		return false, fmt.Errorf("未选择保存路径")
	}

	err = os.WriteFile(selection, []byte(a.content), 0644)
	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "保存文件错误",
			Message: err.Error(),
		})
		return false, err
	}

	a.currentFile = selection

	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   "保存成功",
		Message: fmt.Sprintf("文件已保存到: %s", selection),
	})

	return true, nil
}

// SetContent 设置当前内容（供前端调用）
func (a *App) SetContent(content string) {
	a.content = content
}

// GetCurrentFile 获取当前文件路径
func (a *App) GetCurrentFile() string {
	return a.currentFile
}

// GetContent 获取当前内容（供前端调用）
func (a *App) GetContent() string {
	return a.content
}

// NewFile 创建新文件
func (a *App) NewFile() {
	a.currentFile = ""
	a.content = ""
	runtime.EventsEmit(a.ctx, "openFile", a.content)
}

func (a *App) Close() {
	os.Exit(0)
}
