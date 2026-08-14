<script setup lang="ts">
import { useSvgIcon } from '@/hooks/use-svg-icon-render'
import useTabStore from '@/store/modules/tab'
import { DropdownInstance } from 'element-plus'
import { computed, ref, watch } from 'vue'

defineOptions({ name: 'ContextMenu' })

const { SvgIconVNode } = useSvgIcon()

const { tabName, isAffix, top, left } = defineProps({
    isAffix: {
        type: Boolean,
        default: false
    },
    top: {
        type: Number,
        default: 0
    },
    left: {
        type: Number,
        default: 0
    },
    tabName: {
        type: String,
        default: ''
    }
})

const visible = defineModel({ type: Boolean, default: false })

const contextInstruction = computed(() => {
    return [
        {
            key: 'closeCurrent',
            label: '关闭标签',
            icon: SvgIconVNode({ icon: 'ant-design:close-outlined', fontSize: 20 }),
            disabled: isAffix
        },
        {
            key: 'closeLeft',
            label: '关闭左侧标签',
            icon: SvgIconVNode({ icon: 'ant-design:vertical-right-outlined', fontSize: 20 }),
            disabled: false
        },
        {
            key: 'closeRight',
            label: '关闭右侧标签',
            icon: SvgIconVNode({ icon: 'ant-design:vertical-left-outlined', fontSize: 20 }),
            disabled: false
        },
        {
            key: 'closeOther',
            label: '关闭其他标签',
            icon: SvgIconVNode({ icon: 'ant-design:column-width-outlined', fontSize: 20 }),
            disabled: false
        },
        {
            key: 'closeAll',
            label: '关闭所有标签',
            icon: SvgIconVNode({ icon: 'ant-design:line-outlined', fontSize: 20 }),
            disabled: false
        }
    ]
})

const { removeTab, clearTabs, clearLeftTabs, clearRightTabs } = useTabStore()

const dropdownRef = ref<DropdownInstance>()

watch(visible, val => {
    if (val) {
        dropdownRef.value!.handleOpen()
    } else {
        dropdownRef.value!.handleClose()
    }
})

function hideDropdown() {
    visible.value = false
}

const dropdownAction: Record<App.Global.DropdownKey, () => void> = {
    closeCurrent() {
        removeTab(tabName)
    },
    closeOther() {
        clearTabs([tabName])
    },
    closeLeft() {
        clearLeftTabs(tabName)
    },
    closeRight() {
        clearRightTabs(tabName)
    },
    closeAll() {
        clearTabs()
    }
}
function handleDropdown(optionKey: App.Global.DropdownKey) {
    dropdownAction[optionKey]?.()
    hideDropdown()
}
</script>

<template>
    <div class="absolute" :style="{ top: `${top}px`, left: `${left}px` }">
        <ElDropdown
            ref="dropdownRef"
            :show-arrow="false"
            trigger="click"
            :teleported="false"
            @command="handleDropdown"
            @visible-change="hideDropdown">
            <span></span>
            <template #dropdown>
                <ElDropdownMenu>
                    <ElDropdownItem
                        v-for="{ label, key, icon, disabled } in contextInstruction"
                        :key="key"
                        :icon="icon"
                        :command="key"
                        :disabled="disabled">
                        {{ label }}
                    </ElDropdownItem>
                </ElDropdownMenu>
            </template>
        </ElDropdown>
    </div>
</template>

<style lang="scss" scoped>
.arrow-hide {
    .el-popper__arrow {
        display: none;
    }
}

.contextmenu {
    list-style-type: none;
    padding: 5px 0;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 400;
    box-shadow: 2px 2px 3px 0 rgba(0, 0, 0, 0.3);

    .contextmenu__item {
        margin: 0;
        padding: 7px 16px;
        cursor: pointer;
        display: flex;
        align-items: center;
        white-space: nowrap;
        line-height: 22px;

        i {
            font-size: 14px;
            margin-right: 6px;
        }
    }
}
</style>
