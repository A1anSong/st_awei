import {createRouter, createWebHistory} from "vue-router";
import Home from "@/pages/Home";
import Admin from "@/pages/admin/Admin";
import Dashboard from "@/pages/admin/Dashboard";
import Login from "@/pages/admin/Login";
import NoRoute from "@/pages/NoRoute";

const routes = [
    {
        path: '/',
        component: Home,
        meta: {
            title: '首页'
        },
    },
    {
        path: '/admin',
        component: Admin,
        meta: {
            title: '面板'
        },
        children: [
            {
                path: '',
                component: Dashboard
            },
        ],
    },
    {
        path: '/admin/login',
        component: Login,
        meta: {
            title: '登录'
        },
    },
    {
        path: '/:pathMatch(.*)*',
        component: NoRoute,
        meta: {
            title: '404 资源不存在'
        },
    },
]

const router = createRouter({
    history: createWebHistory(),
    routes: routes,
})

router.beforeEach((to, from, next) => {
    if (to.meta.title) {
        document.title = to.meta.title
    }
    next()
})

export default router