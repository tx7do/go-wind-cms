<script lang="ts" setup>
import { computed, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { $t } from '@vben/locales';
import { preferences } from '@vben/preferences';
import { useUserStore } from '@vben/stores';

import { Col, notification, Row, Upload } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import {
  fetchUserProfile,
  genderList,
  getMe,
  makeUpdateMask,
  updateMyUserInfo,
  uploadMediaAsset,
  type identityservicev1_User as User,
} from '#/api';

const data = ref<null | User>();
const submitLoading = ref(false);
const userStore = useUserStore();

const avatar = computed(() => {
  return data.value?.avatar ?? preferences.app.defaultAvatar;
});

async function refreshProfile() {
  data.value = await getMe();
  await baseFormApi.setValues(data.value || {});
  // 同步全局 userStore，使布局等处的头像展示随之更新。
  const fresh = await fetchUserProfile();
  if (fresh) {
    userStore.setUserInfo(fresh as any);
  }
}

async function handleUploadAvatar(option: any) {
  const { file, onSuccess, onError } = option;
  try {
    const resp = await uploadMediaAsset({}, file as File);
    const url = (resp as any).objectName || '';
    if (!url) {
      throw new Error('avatar upload returned empty url');
    }
    await updateMyUserInfo({
      id: data.value?.id,
      data: { avatar: url } as any,
      updateMask: makeUpdateMask(['avatar']),
    } as any);
    notification.success({ message: $t('ui.notification.update_success') });
    await refreshProfile();
    onSuccess?.({}, file);
  } catch (error) {
    console.error('avatar upload failed:', error);
    notification.error({ message: $t('ui.notification.update_failed') });
    onError?.(error as any);
  }
  return false;
}

const [BaseForm, baseFormApi] = useVbenForm({
  showDefaultActions: false,
  commonConfig: {
    componentProps: {
      class: 'w-full',
    },
  },
  schema: [
    {
      fieldName: 'nickname',
      component: 'Input',
      label: $t('page.user.table.nickname'),
    },
    {
      fieldName: 'realname',
      component: 'Input',
      label: $t('page.user.table.realname'),
    },
    {
      fieldName: 'email',
      component: 'Input',
      label: $t('page.user.table.email'),
    },
    {
      fieldName: 'mobile',
      component: 'Input',
      label: $t('page.user.table.mobile'),
    },
    {
      fieldName: 'telephone',
      component: 'Input',
      label: $t('page.user.table.telephone'),
    },
    {
      fieldName: 'gender',
      component: 'Select',
      label: $t('page.user.table.gender'),
      componentProps: {
        filterOption: (input: string, option: any) =>
          option.label.toLowerCase().includes(input.toLowerCase()),
        allowClear: true,
        showSearch: true,
        options: genderList,
        placeholder: $t('ui.placeholder.select'),
      },
    },
    {
      fieldName: 'region',
      component: 'Input',
      label: $t('page.user.table.region'),
    },
    {
      fieldName: 'address',
      component: 'Input',
      label: $t('page.user.table.address'),
    },
    {
      fieldName: 'description',
      component: 'Textarea',
      label: $t('page.user.table.description'),
    },
  ],
});

async function handleSubmit() {
  console.log('submit');

  // 校验输入的数据
  const validate = await baseFormApi.validate();
  if (!validate.valid) {
    return;
  }

  setLoading(true);

  // 获取表单数据
  const values = await baseFormApi.getValues();

  try {
    await updateMyUserInfo({
      id: data.value?.id,
      data: { ...values } as any,
      updateMask: makeUpdateMask(Object.keys(values)),
    });

    notification.success({
      message: $t('ui.notification.update_success'),
    });
  } catch {
    notification.error({
      message: $t('ui.notification.update_failed'),
    });
  } finally {
    setLoading(false);
  }
}

function setLoading(loading: boolean) {
  submitLoading.value = loading;
}

/**
 * 重新加载用户信息
 */
async function reload() {
  await refreshProfile();
}

reload();
</script>

<template>
  <Page
    :title="$t('page.user.profile.tab.basicSettings')"
    :body-style="{ padding: 0 }"
    class="edge-card"
    style="margin: 0"
  >
    <Row :gutter="24">
      <Col :span="14">
        <BaseForm />
      </Col>
      <Col :span="10">
        <div class="change-avatar">
          <div class="mb-2">{{ $t('page.user.table.avatar') }}</div>
          <img :src="avatar" :alt="$t('page.user.table.avatar')" />
          <Upload
            :custom-request="handleUploadAvatar"
            :show-upload-list="false"
            accept="image/*"
            :before-upload="() => false"
          >
            <a-button>{{ $t('page.user.button.changeAvatar') }}</a-button>
          </Upload>
        </div>
      </Col>
    </Row>
    <a-button
      type="primary"
      :loading="submitLoading"
      @click="handleSubmit"
    >
      {{ $t('page.user.button.updateUserInfo') }}
    </a-button>
  </Page>
</template>

<style lang="less" scoped>
.change-avatar {
  img {
    display: block;
    margin-bottom: 15px;
    border-radius: 50%;
  }
}

.edge-card {
  .ant-card-body {
    padding: 0 !important;
  }
}
</style>
