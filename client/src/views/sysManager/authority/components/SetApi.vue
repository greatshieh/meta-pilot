<script setup lang="ts">
import { apiApi, casbinApi } from '@/api'
import { TreeInstance, TreeNode, TreeNodeData } from 'element-plus'
import { ref, shallowRef, watch } from 'vue'

defineOptions({ name: 'SetApi' })

const { authorityId } = defineProps({
    maxHeight: {
        type: String,
        default: '400px'
    },
    authorityId: {
        type: Number,
        default: 0
    }
})

interface ApiGroup {
    ID: string
    description: string
    children: Api.ApiMetaData[]
}

const apiList = ref<ApiGroup[]>()
const apiTreeRef = shallowRef<TreeInstance>()
const selectedApiIds = ref<string[]>([])

const filterDesc = ref('')
const filterPath = ref('')

const apiTreePorps = {
    label: (data: TreeNodeData, _node: TreeNode) => data.description,
    disabled: (data: TreeNodeData, _node: TreeNode) => data.required,
    children: 'children'
}

function transformArrytoApiTree(apis: Api.ApiMetaData[]) {
    const apiObj: Record<string, (Api.ApiMetaData & { onlyId: string })[]> = {}
    apis.forEach(item => {
        if (Object.prototype.hasOwnProperty.call(apiObj, item.apiGroup)) {
            apiObj[item.apiGroup].push({ ...item, onlyId: `${item.path}-${item.method}` })
        } else {
            Object.assign(apiObj, { [item.apiGroup]: [{ ...item, onlyId: `${item.path}-${item.method}` }] })
        }
    })
    const apiTree = []
    for (const key in apiObj) {
        const treeNode = {
            ID: key,
            description: key + '组',
            children: apiObj[key]
        }
        apiTree.push(treeNode)
    }
    return apiTree
}

function init() {
    casbinApi
        .getPolicyPathByAuthorityId(authorityId)
        .then(res => {
            if (res) {
                res.forEach(item => {
                    selectedApiIds.value.push(`${item.path}-${item.method}`)
                })
            }
            // 获取所有api
            return apiApi.getAllApis()
        })
        .then(res => {
            apiList.value = transformArrytoApiTree(res)
            apiList.value?.forEach(item => {
                item.children.forEach(child => {
                    if (!child.required) return
                    if (selectedApiIds.value.indexOf(`${child.path}-${child.method}`) === -1) {
                        selectedApiIds.value.push(`${child.path}-${child.method}`)
                    }
                })
            })
        })
}

init()

function setApiAuthority() {
    const checkArr = (apiTreeRef.value?.getCheckedNodes(true, true) || []) as Api.ApiMetaData[]
    if (checkArr.length === 0) {
        ElMessage.error('请选择api')
        return
    }

    var casbinInfos: Api.Cabin.PolicyPath[] = []
    checkArr.forEach(item => {
        var casbinInfo = {
            path: item.path,
            method: item.method
        }
        casbinInfos.push(casbinInfo)
    })

    casbinApi.setPolicy(authorityId, { casbinInfos }).then(() => {
        ElMessage.success('api 权限设置成功')
    })
}

function filterNode(value: [string, string], data: TreeNodeData) {
    if (data.children?.length > 0) return false

    if (value[0] === '' && value[1] === '') return true

    if (value[0] !== '' && data.description?.indexOf(value[0]) !== -1) return true
    if (value[1] !== '' && data.path?.indexOf(value[1]) !== -1) return true

    return false
}

watch([filterDesc, filterPath], val => {
    apiTreeRef.value?.filter(val)
})
</script>

<template>
    <div class="flex items-center justify-between gap-4 mb-2">
        <div class="flex items-center gap-2 flex-1">
            <ElInput v-model="filterDesc" size="small" placeholder="api 描述"></ElInput>
            <ElInput v-model="filterPath" size="small" placeholder="api 路径"></ElInput>
        </div>
        <ElButton type="primary" size="small" @click="setApiAuthority">确定</ElButton>
    </div>
    <ElScrollbar max-height="400px">
        <ElTree
            ref="apiTreeRef"
            :data="apiList"
            :props="apiTreePorps"
            node-key="onlyId"
            :default-checked-keys="selectedApiIds"
            show-checkbox
            :filter-node-method="filterNode"
            default-expand-all>
            <template #default="{ _node, data }">
                <div class="flex items-center justify-between flex-1 pr-4" :class="{ 'c-primary': data.required }">
                    <span>{{ data.description }} {{ data.required ? '(必选)' : '' }}</span>
                    <span v-if="data.method && data.path">({{ data.method }}) {{ data.path }}</span>
                </div>
            </template>
        </ElTree>
    </ElScrollbar>
</template>
