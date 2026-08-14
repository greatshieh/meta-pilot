<template>
    <el-dropdown trigger="click" class="cursor-pointer">
        <div class="profile-container">
            <img :src="adminStore.userInfo.avatar" class="avatar-container" />
            <span>{{ adminStore.userInfo.nickName }}</span>
        </div>
        <template #dropdown>
            <el-dropdown-menu>
                <el-dropdown-item>当前角色: {{ adminStore.userInfo.authority.authorityName }}</el-dropdown-item>
                <el-dropdown-item
                    v-for="item in adminStore.userInfo.authorities.filter(i => i.authorityId !== adminStore.userInfo.authority.authorityId)"
                    :key="item.authorityId"
                    @click="changeUserAuth(item.authorityId)">
                    切换为: {{ item.authorityName }}
                </el-dropdown-item>
                <el-dropdown-item @click="gotoProfile">
                    <SvgIcon icon="mdi:face-man-profile"></SvgIcon>
                    <span class="ml-2">个人中心</span>
                </el-dropdown-item>
                <el-dropdown-item divided @click="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
        </template>
    </el-dropdown>
</template>

<script lang="ts" setup name="UserProfile">
import { useRoute, useRouter } from 'vue-router'
import useAdminStore from '@/store/modules/admin'
import { sysUserApi } from '@/api'

defineOptions({ name: 'UserProfile' })

const adminStore = useAdminStore()

const $route = useRoute()
const $router = useRouter()

const gotoProfile = () => {
    $router.push({ name: 'profile' })
}

const logout = async () => {
    adminStore.logout()
    $router.push(`/login?redirect=${$route.fullPath}`)
}

function changeUserAuth(id: number) {
    sysUserApi.setUserAuthority({ authorityId: id }).then(() => {
        window.location.reload()
    })
}
</script>

<style lang="scss" scoped>
.profile-container {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    padding: 0 13px;

    .avatar-container {
        border-radius: 9999px;
        height: 24px;
        width: 24px;
        margin-right: 10px;
    }
}
</style>
