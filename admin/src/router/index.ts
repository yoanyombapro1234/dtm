import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { IMenubarList } from '../type/store/layout';
import { components } from './asyncRouter';

const Components: IObject<() => Promise<typeof import('*.vue')>> = Object.assign({}, components, {
  LayoutHeader: (() => import('/@/layout/index.vue')) as unknown as () => Promise<typeof import('*.vue')>,
  LayoutMain: (() => import('/@/layout/aside.vue')) as unknown as () => Promise<typeof import('*.vue')>
})

export const allowRouter: Array<IMenubarList> = [
  {
    name: 'Admin',
    path: '/',
    redirect: '/admin/global-transactions/all',
    component: Components['LayoutHeader'],
    meta: { title: 'Admin', activeMenu: '/admin' },
    children: [
      {
        name: 'Nodes',
        path: '/admin/nodes',
        component: Components['LayoutMain'],
        meta: { title: 'Nodes' },
        children: [
          {
            name: 'LivingNodes',
            path: '/admin/nodes/living',
            component: Components['LivingNodes'],
            meta: { title: 'Living Nodes' },
          }
        ]
      }, {
        name: 'GlobalTransactions',
        path: '/admin/global-transactions',
        component: Components['LayoutMain'],
        meta: { title: 'Global Transactions' },
        children: [
          {
            name: 'AllTransactions',
            path: '/admin/global-transactions/all',
            component: Components['AllTransactions'],
            meta: { title: 'All Transactions' },
          }, {
            name: 'UnfinishedTransactions',
            path: '/admin/global-transactions/unfinished',
            component: Components['UnfinishedTransactions'],
            meta: { title: 'Unfinished Transactions' },
          }
        ]
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes: allowRouter as RouteRecordRaw[]
})

export default router
