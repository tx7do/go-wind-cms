<script setup lang="ts">
import { ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { notification } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import { apiClient } from '#/api';

const data = ref();

const [BaseForm, baseFormApi] = useVbenForm({
  showDefaultActions: false,
  // 所有表单项共用，可单独在表单内覆盖
  commonConfig: {
    // 所有表单项
    componentProps: {
      class: 'w-full',
    },
  },
  schema: [
    {
      component: 'Textarea',
      fieldName: 'summary',
      label: $t('page.post.summary'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
      },
      rules: 'required',
    },
  ],
});

const [Modal, modalApi] = useVbenModal({
  onCancel() {
    modalApi.close();
  },

  async onConfirm() {
    console.log('onConfirm');

    // 校验输入的数据
    const validate = await baseFormApi.validate();
    if (!validate.valid) {
      return;
    }

    setLoading(true);

    // 获取表单数据
    const values = await baseFormApi.getValues();

    if (values.new_password !== values.confirm_password) {
      notification.error({
        message: $t('page.notification.password_mismatch'),
      });

      setLoading(false);
      modalApi.close();
      return;
    }

    try {
      await apiClient.postService.Create({ data: values });

      setLoading(false);
      modalApi.close();

      notification.success({
        message: $t('ui.notification.create_success'),
      });
    } catch {
      setLoading(false);

      notification.error({
        message: $t('ui.notification.create_failed'),
      });
    }
  },

  onOpenChange(isOpen: boolean) {
    if (isOpen) {
      // 获取传入的数据
      data.value = modalApi.getData<any>();

      setLoading(false);

      console.log('onOpenChange', data.value);
    }
  },
});

function setLoading(loading: boolean) {
  modalApi.setState({ confirmLoading: loading });
}
</script>

<template>
  <Modal :title="$t('page.post.button.publish')">
    <BaseForm />
  </Modal>
</template>

<style scoped></style>
