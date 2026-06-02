import { defineNuxtPlugin } from '#app';
import { initPreferences } from '@/core/preferences';

export default defineNuxtPlugin(async () => {
  if (import.meta.server) return;
  await initPreferences({
    namespace: 'gowind-cms',
  });
});
