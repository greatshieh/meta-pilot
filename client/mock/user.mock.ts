import { defineMock, MockHttpItem } from "vite-plugin-mock-dev-server";

const requests: MockHttpItem[] = [
    {
        url: '/api/v1/users/login',
        method: 'POST',
        body: () => {
            return {
                "code": 20000,
                "data": {
                    "user": {
                        "id": 1,
                        "CreatedAt": "2026-07-10T23:51:05.625+08:00",
                        "UpdatedAt": "2026-07-11T21:42:50.209+08:00",
                        "DeletedAt": null,
                        "uuid": "f03e4c90-8421-40d3-b385-c8b8e08d2679",
                        "userName": "admin",
                        "nickName": "超级管理员",
                        "avatar": "https://img.tuxiangyan.com/zb_users/upload/2023/02/202302091675904134770942.jpg",
                        "introduction": "",
                        "authorityID": 200,
                        "authority": {
                            "createdAt": "0001-01-01T00:00:00Z",
                            "updatedAt": "0001-01-01T00:00:00Z",
                            "DeletedAt": null,
                            "authorityID": 0,
                            "authorityName": "",
                            "parentID": null,
                            "authorized_sub_roles": null,
                            "children": null,
                            "menus": null,
                            "users": null,
                            "defaultRouter": ""
                        },
                        "authorities": null,
                        "phone": "17700000000",
                        "email": "88888888@qq.com",
                        "isStaff": true,
                        "isActive": true,
                        "lastLogin": "2026-07-11 21:42:50"
                    },
                    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVVUlEIjoiZjAzZTRjOTAtODQyMS00MGQzLWIzODUtYzhiOGUwOGQyNjc5IiwiSUQiOjEsIlVzZXJOYW1lIjoiYWRtaW4iLCJOaWNrTmFtZSI6Iui2hee6p-euoeeQhuWRmCIsIkF1dGhvcml0eUlkIjoyMDAsIlJvbGUiOiIiLCJCdWZmZXJUaW1lIjo4NjQwMCwiaXNzIjoibWV0YS1waWxvdCIsImV4cCI6MTc4NDU1NjYyMCwibmJmIjoxNzgzOTUwODIwfQ.dC_aZz3S483NbP64rXxtAOGHPitWwqE_8Vl5NsbgsFg",
                    "expiresAt": 1784556620000
                },
                "message": "登录成功"
            }
        },
    },
    {
        url: "/api/v1/users/me",
        method: "GET",
        body: {
            "code": 20000,
            "data": {
                "id": 1,
                "CreatedAt": "2026-07-10T23:51:05.625+08:00",
                "UpdatedAt": "2026-07-13T22:10:20.119+08:00",
                "DeletedAt": null,
                "uuid": "f03e4c90-8421-40d3-b385-c8b8e08d2679",
                "userName": "admin",
                "nickName": "超级管理员",
                "avatar": "https://img.tuxiangyan.com/zb_users/upload/2023/02/202302091675904134770942.jpg",
                "introduction": "",
                "authorityID": 200,
                "authority": {
                    "createdAt": "2026-07-10T23:51:05.531+08:00",
                    "updatedAt": "2026-07-10T23:51:05.544+08:00",
                    "DeletedAt": null,
                    "authorityID": 200,
                    "authorityName": "系统管理员",
                    "parentID": 0,
                    "authorized_sub_roles": null,
                    "children": null,
                    "menus": null,
                    "users": null,
                    "defaultRouter": "dashboard"
                },
                "authorities": [
                    {
                        "createdAt": "2026-07-10T23:51:05.531+08:00",
                        "updatedAt": "2026-07-10T23:51:05.544+08:00",
                        "DeletedAt": null,
                        "authorityID": 200,
                        "authorityName": "系统管理员",
                        "parentID": 0,
                        "authorized_sub_roles": null,
                        "children": null,
                        "menus": null,
                        "users": null,
                        "defaultRouter": "dashboard"
                    }
                ],
                "phone": "17700000000",
                "email": "88888888@qq.com",
                "isStaff": true,
                "isActive": true,
                "lastLogin": "2026-07-13 22:10:20"
            },
            "message": "SUCCESS"
        }
    }]

export default defineMock(requests)
