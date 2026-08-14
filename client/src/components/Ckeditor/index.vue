<template>
    <ElDialog
        v-if="shouldShowDialog"
        v-model="dialogShow"
        :title="title"
        destroy-on-close
        :close-on-click-modal="false"
        :close-on-press-escape="false"
        :show-close="false"
        append-to-body
        @close="dialogClose">
        <div id="ckeditor"></div>

        <template #footer>
            <ElButton type="primary" @click="handleUpdate">确定</ElButton>
            <ElButton type="info" @click="handleCancel">取消</ElButton>
        </template>
    </ElDialog>
</template>

<script setup lang="ts" name="Editor">
import MathType from '@wiris/mathtype-ckeditor5/dist/index.js'
import {
    AutoImage,
    Autosave,
    Bold,
    ClassicEditor,
    type EditorConfig,
    type FileLoader,
    FontSize,
    Heading,
    ImageBlock,
    ImageInline,
    ImageInsertViaUrl,
    ImageResize,
    ImageStyle,
    ImageToolbar,
    ImageUpload,
    Italic,
    RemoveFormat,
    SelectAll,
    SpecialCharacters,
    SpecialCharactersArrows,
    SpecialCharactersCurrency,
    SpecialCharactersEssentials,
    SpecialCharactersLatin,
    SpecialCharactersMathematical,
    SpecialCharactersText,
    Subscript,
    Superscript,
    Table,
    TableToolbar,
    Underline,
    Undo
} from 'ckeditor5'
import translations from 'ckeditor5/translations/zh-cn.js'
import 'ckeditor5/ckeditor5.css'
import { ElMessage, ElMessageBox } from 'element-plus'
import { computed, nextTick, ref, toRaw, watch } from 'vue'
import { UploadAdapter } from './UploadImageAdapter'

const props = defineProps({
    editorContent: {
        type: String,
        default: ''
    },
    title: {
        type: String,
        default: '试题编辑器'
    },
    editorShow: {
        type: Boolean,
        defaule: false
    }
})

const $emits = defineEmits(['update', 'cancel', 'update:editorShow'])

const editorInstance = ref()

const dialogShow = computed({
    get: () => props.editorShow,
    set: val => $emits('update:editorShow', val)
})

const shouldShowDialog = ref(false)

const config: EditorConfig = {
    toolbar: {
        items: [
            'undo',
            'redo',
            '|',
            'fontSize',
            '|',
            'bold',
            'italic',
            'underline',
            'subscript',
            'superscript',
            'removeFormat',
            '|',
            'specialCharacters',
            'insertTable',
            '|',
            'insertImage',
            '|',
            'MathType'
        ],
        shouldNotGroupWhenFull: false
    },
    plugins: [
        AutoImage,
        Autosave,
        Bold,
        FontSize,
        Heading,
        ImageBlock,
        ImageInline,
        ImageInsertViaUrl,
        ImageResize,
        ImageStyle,
        ImageToolbar,
        ImageUpload,
        Italic,
        RemoveFormat,
        SelectAll,
        SpecialCharacters,
        SpecialCharactersArrows,
        SpecialCharactersCurrency,
        SpecialCharactersEssentials,
        SpecialCharactersLatin,
        SpecialCharactersMathematical,
        SpecialCharactersText,
        Subscript,
        Superscript,
        Table,
        TableToolbar,
        Underline,
        Undo,
        MathType
    ],
    fontSize: {
        options: [10, 12, 14, 'default', 18, 20, 22],
        supportAllValues: true
    },
    image: {
        toolbar: ['imageTextAlternative', '|', 'imageStyle:inline', 'imageStyle:wrapText', 'imageStyle:breakText', '|', 'resizeImage']
    },
    language: 'zh-cn',
    placeholder: '在此处输入内容，支持格式调整',
    table: {
        contentToolbar: ['tableColumn', 'tableRow', 'mergeTableCells']
    },
    translations: [translations]
}

watch(
    () => dialogShow.value,
    val => {
        if (val) {
            shouldShowDialog.value = true
            nextTick(() => {
                ClassicEditor.create(document.querySelector('#ckeditor') as HTMLElement, config).then((editor: ClassicEditor) => {
                    editorInstance.value = editor
                    // 设置初始化内容
                    if (props.editorContent !== '') {
                        editor.setData(props.editorContent)
                    }

                    // 加载上传插件
                    editor.plugins.get('FileRepository').createUploadAdapter = (loader: FileLoader) => {
                        return new UploadAdapter(loader)
                    }
                })
            })
        }
    },
    { immediate: false }
)

function handleUpdate() {
    const editor = toRaw(editorInstance.value)
    $emits('update', editor.getData())
}

function handleCancel() {
    ElMessageBox.confirm('试题未保存，是否继续? 如果离开，系统不会保存修改的内容。', '提示', {
        cancelButtonText: '取消',
        confirmButtonText: '确认退出',
        type: 'warning'
    }).then(() => {
        ElMessage({
            type: 'warning',
            message: '编辑未保存'
        })
        $emits('cancel')
    })
}

function dialogClose() {
    const editor = toRaw(editorInstance.value)

    document.querySelectorAll('.wrs_modal_dialogContainer')?.forEach(item => {
        document.body.removeChild(item)
    })
    editor.destroy()
    shouldShowDialog.value = false
}
</script>

<style lang="scss">
.ck-content {
    line-height: 1.6;
    word-break: break-word;
    min-height: 300px;
}

// .el-dialog__body {
//   margin-top: 150px;
// }
</style>
