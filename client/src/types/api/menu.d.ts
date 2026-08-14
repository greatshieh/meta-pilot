declare namespace Api.Menu {
    interface MenuMeta {
        id: number
        name: string
        path: string
        redirect?: string
        component: string
        menuId?: number
        parentId?: number
        sort?: number
        children: MenuMeta[]
        meta: RouteMeta
    }
}
