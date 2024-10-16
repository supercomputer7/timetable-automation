<template>
  <v-form
    ref="form"
    lazy-validation
    v-model="formValid"
    v-if="isLoaded()"
  >
    <v-alert
      :value="formError"
      color="error"
      class="pb-2"
    >{{ formError }}
    </v-alert>

    <v-alert
      color="blue"
      dark
      icon="mdi-source-fork"
      dismissible
      v-model="commitAvailable"
      prominent
    >
      <div
        style="font-weight: bold;"
      >{{ (item.commit_hash || '').substr(0, 10) }}
      </div>
      <div v-if="sourceTask && sourceTask.commit_message">{{ sourceTask.commit_message }}</div>
    </v-alert>

    <v-select
      v-if="template.type === 'deploy'"
      v-model="item.build_task_id"
      :label="$t('buildVersion')"
      :items="buildTasks"
      item-value="id"
      :item-text="(itm) => getTaskMessage(itm)"
      :rules="[v => !!v || $t('build_version_required')]"
      required
      :disabled="formSaving"
    />

    <v-text-field
      v-model="item.message"
      :label="$t('messageOptional')"
      :disabled="formSaving"
    />

    <TaskParamsForm v-if="template.app === 'ansible'" v-model="item" :app="template.app" />
    <TaskParamsForm v-else v-model="item.params" :app="template.app" />

    <ArgsPicker
      v-if="template.allow_override_args_in_task"
      :vars="args"
      @change="setArgs"
    />

  </v-form>
</template>
<script>
/* eslint-disable import/no-extraneous-dependencies,import/extensions */

import ItemFormBase from '@/components/ItemFormBase';
import axios from 'axios';
// import { codemirror } from 'vue-codemirror';
import 'codemirror/lib/codemirror.css';
import 'codemirror/mode/vue/vue.js';
import 'codemirror/addon/lint/json-lint.js';
import 'codemirror/addon/display/placeholder.js';
import TaskParamsForm from '@/components/TaskParamsForm.vue';
import ArgsPicker from '@/components/ArgsPicker.vue';

export default {
  mixins: [ItemFormBase],
  props: {
    templateId: Number,
    sourceTask: Object,
  },
  components: {
    ArgsPicker,
    TaskParamsForm,
    // codemirror,
  },
  data() {
    return {
      template: null,
      buildTasks: null,
      commitAvailable: null,
      cmOptions: {
        tabSize: 2,
        mode: 'application/json',
        lineNumbers: true,
        line: true,
        lint: true,
        indentWithTabs: false,
      },
      // advancedOptions: false,
    };
  },
  computed: {
    args() {
      return JSON.parse(this.item.arguments || '[]');
    },
  },

  watch: {
    needReset(val) {
      if (val) {
        this.item.template_id = this.templateId;
      }
    },

    templateId(val) {
      this.item.template_id = val;
    },

    sourceTask(val) {
      this.assignItem(val);
    },

    commitAvailable(val) {
      if (val == null) {
        this.commit_hash = null;
      }
    },
  },

  methods: {
    setArgs(args) {
      this.item.arguments = JSON.stringify(args || []);
    },

    getTaskMessage(task) {
      let buildTask = task;

      while (buildTask.version == null && buildTask.build_task != null) {
        buildTask = buildTask.build_task;
      }

      if (!buildTask) {
        return '';
      }

      return buildTask.version + (buildTask.message ? ` — ${buildTask.message}` : '');
    },

    assignItem(val) {
      const v = val || {};

      if (this.item == null) {
        this.item = {};
      }

      Object.keys(v).forEach((field) => {
        this.item[field] = v[field];
      });

      this.commitAvailable = v.commit_hash != null;
    },

    isLoaded() {
      return this.item != null
        && this.template != null
        && this.buildTasks != null;
    },

    async afterLoadData() {
      this.assignItem(this.sourceTask);

      this.item.template_id = this.templateId;

      if (!this.item.params) {
        this.item.params = {};
      }

      // this.advancedOptions = this.item.arguments != null;

      this.template = (await axios({
        keys: 'get',
        url: `/api/project/${this.projectId}/templates/${this.templateId}`,
        responseType: 'json',
      })).data;

      this.buildTasks = this.template.type === 'deploy' ? (await axios({
        keys: 'get',
        url: `/api/project/${this.projectId}/templates/${this.template.build_template_id}/tasks?status=success`,
        responseType: 'json',
      })).data.filter((task) => task.status === 'success') : [];

      if (this.item.build_task_id == null
        && this.buildTasks.length > 0
        && this.buildTasks.length > 0) {
        this.item.build_task_id = this.buildTasks[0].id;
      }
    },

    getItemsUrl() {
      return `/api/project/${this.projectId}/tasks`;
    },
  },
};
</script>
