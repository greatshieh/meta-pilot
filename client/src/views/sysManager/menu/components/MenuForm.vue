<script setup lang="ts">
import { cloneDeep, isEqual } from 'lodash'
import { computed, PropType, ref } from 'vue'
import useAppStore from '@/store/modules/app'
import useRouteStore from '@/store/modules/route'
import { invalidateCache, Method } from 'alova'
import { menuApi } from '@/api'

defineOptions({ name: 'MenuForm' })

const props = defineProps({
    title: {
        type: String,
        default: '编辑菜单'
    },
    show: {
        type: Boolean,
        default: false
    },
    menu: {
        type: Object as PropType<Api.Menu.MenuMeta>,
        required: true,
        default: () => ({})
    }
})

const $emits = defineEmits(['update:show'])

const appStore = useAppStore()
const routeStore = useRouteStore()

const editMenu = ref(cloneDeep(props.menu))

const visible = computed({
    get: () => props.show,
    set: val => $emits('update:show', val)
})

async function handleSubmit() {
    if (!isEqual(editMenu.value, props.menu)) {
        if (props.title.includes('编辑')) {
            await menuApi.updateMenu(editMenu.value)
        } else {
            await menuApi.addMenu(editMenu.value)
        }
        // 重新获取异步路由
        invalidateCache(menuApi.getAsyncMenu() as Method)
        await routeStore.initAuthRoute()
        // 刷新页面
        appStore.reloadPage()
    }
    visible.value = false
}
</script>

<template>
    <ElDialog v-model="visible" :title="title" :show-close="false">
        <template #header>
            <div class="flex items-center justify-between">
                <span>{{ title }}</span>
                <div>
                    <ElButton type="info" size="small" @click="visible = false">取消</ElButton>
                    <ElButton type="primary" size="small" @click="handleSubmit">提交</ElButton>
                </div>
            </div>
        </template>
        <ElForm :model="editMenu" label-position="top" size="small">
            <div class="flex items-center justify-start gap-8">
                <ElFormItem label="展示名称" required>
                    <ElInput v-model="editMenu.meta.title" style="width: 120px"></ElInput>
                </ElFormItem>
                <ElFormItem label="文件路径" required class="flex-1">
                    <ElInput v-model="editMenu.component"></ElInput>
                </ElFormItem>
            </div>
            <div class="flex items-center justify-start gap-8">
                <ElFormItem label="路由name" required>
                    <ElInput v-model="editMenu.name" style="width: 120px"></ElInput>
                </ElFormItem>
                <ElFormItem label="路由path" required class="w-40%">
                    <ElInput v-model="editMenu.path"></ElInput>
                </ElFormItem>
                <ElFormItem label="重定向path" class="flex-1">
                    <ElInput v-model="editMenu.redirect"></ElInput>
                </ElFormItem>
            </div>

            <div class="flex items-center justify-start gap-8">
                <ElFormItem label="父菜单ID" required>
                    <ElInput v-model="editMenu.parentId" disabled></ElInput>
                </ElFormItem>
                <ElFormItem label="排序">
                    <ElInputNumber v-model="editMenu.sort"></ElInputNumber>
                </ElFormItem>

                <ElFormItem label="图标">
                    <ElInput v-model="editMenu.meta.icon" style="width: 200px">
                        <template v-if="editMenu.meta.icon" #suffix>
                            <ElIcon size="18">
                                <SvgIcon :icon="editMenu.meta.icon" />
                            </ElIcon>
                        </template>
                    </ElInput>
                </ElFormItem>

                <ElFormItem label="隐藏">
                    <ElSwitch v-model="editMenu.meta.hidden"></ElSwitch>
                </ElFormItem>

                <ElFormItem label="高亮菜单name">
                    <ElInput v-model="editMenu.meta.activeName" style="width: 120px"></ElInput>
                </ElFormItem>

                <ElFormItem label="固定标签页">
                    <ElSwitch v-model="editMenu.meta.affix"></ElSwitch>
                </ElFormItem>
            </div>
        </ElForm>
    </ElDialog>
</template>
