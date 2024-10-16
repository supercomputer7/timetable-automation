import Vue from 'vue';
import VueRouter from 'vue-router';
import RestoreProject from '@/views/project/RestoreProject.vue';
import Tasks from '@/views/Tasks.vue';
import Schedule from '../views/project/Schedule.vue';
import History from '../views/project/History.vue';
import Activity from '../views/project/Activity.vue';
import Settings from '../views/project/Settings.vue';
import Templates from '../views/project/Templates.vue';
import TemplateView from '../views/project/TemplateView.vue';
import Users from '../views/Users.vue';
import Auth from '../views/Auth.vue';
import New from '../views/project/New.vue';
import Apps from '../views/Apps.vue';
import Runners from '../views/Runners.vue';

Vue.use(VueRouter);

const routes = [
  {
    path: '/project/new',
    component: New,
  },
  {
    path: '/project/restore',
    component: RestoreProject,
  },
  {
    path: '/project/:projectId',
    redirect: '/project/:projectId/history',
  },
  {
    path: '/project/:projectId/history',
    component: History,
  },
  {
    path: '/project/:projectId/activity',
    component: Activity,
  },
  {
    path: '/project/:projectId/schedule',
    component: Schedule,
  },
  {
    path: '/project/:projectId/settings',
    component: Settings,
  },
  {
    path: '/project/:projectId/templates',
    component: Templates,
  },
  {
    path: '/project/:projectId/views/:viewId/templates',
    component: Templates,
  },
  {
    path: '/project/:projectId/templates/:templateId',
    component: TemplateView,
  },
  {
    path: '/project/:projectId/views/:viewId/templates/:templateId',
    component: TemplateView,
  },
  {
    path: '/auth/login',
    component: Auth,
  },
  {
    path: '/users',
    component: Users,
  },
  {
    path: '/runners',
    component: Runners,
  },
  {
    path: '/tasks',
    component: Tasks,
  },
  {
    path: '/apps',
    component: Apps,
  },
];

const router = new VueRouter({
  mode: 'history',
  routes,
});

export default router;
