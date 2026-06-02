import { defineNuxtPlugin } from '#app';
import { initPreferences } from '@/core/preferences';

export default defineNuxtPlugin(async () => {
  await initPreferences({
    namespace: 'gowind-cms',
  });
});
