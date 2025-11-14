<template>
    <div ref="editorRef"></div>
</template>

<script>
import { SetContent } from '../wailsjs/go/main/App';
import { EventsOn } from '../wailsjs/runtime';

import Editor from '@toast-ui/editor';
import '@toast-ui/editor/dist/i18n/zh-cn';

//uml插件
import uml from '@toast-ui/editor-plugin-uml';
//代码高亮插件
import codeSyntaxHighlight from '@toast-ui/editor-plugin-code-syntax-highlight';
//图表插件
import chart from '@toast-ui/editor-plugin-chart';
//文本颜色选择插件
import colorSyntax from '@toast-ui/editor-plugin-color-syntax';
//表格合并插件
import tableMergedCell from '@toast-ui/editor-plugin-table-merged-cell';

export default {
    data() {
        return {
            content: '',
            initialEditType: 'wysiwyg',
            editorInstance: null,
            controller: null,
        };
    },
    mounted() {
        this.controller = this.$refs.editorRef;
        EventsOn('openFile', (data) => {
            if (this.editorInstance) {
                try {
                    this.editorInstance.destroy();
                } catch (e) {
                    console.warn('销毁旧编辑器时出错:', e)
                }
                this.editorInstance = null
            }
            this.content = data;
            this.initEdit();
        });
        EventsOn('changeMode', (data) => {
            this.editorInstance.destroy();
            this.editorInstance = null;
            this.initialEditType = data;
            this.initEdit();
        });
        this.initEdit();

    },
    methods: {
        initEdit() {
            this.editorInstance = new Editor({
                el: this.controller,
                height: "95vh",
                previewStyle: 'vertical',
                initialEditType: this.initialEditType,
                initialValue: this.content,
                language: 'zh-CN',
                usageStatistics: false,
                useCommandShortcut: false,
                hideModeSwitch: true,
                plugins: [codeSyntaxHighlight, chart, uml, colorSyntax, tableMergedCell],
                hooks: {
                    addImageBlobHook: async (blob, callback) => {
                        //转换成base64 生成md的图片地址
                        let base64 = await new Promise((resolve, reject) => {
                            const reader = new FileReader();
                            reader.onload = (e) => {
                                resolve(e.target.result);
                            };
                            reader.onerror = (e) => {
                                reject(e);
                            };
                            reader.readAsDataURL(blob);
                        });
                        callback(base64, '');
                    }
                }
            });
            this.editorInstance.on('change', () => {
                let markdown = this.editorInstance.getMarkdown();
                if (markdown !== this.content) {
                    this.content = markdown;
                    SetContent(this.content);
                }
            });
        }
    }
};
</script>
