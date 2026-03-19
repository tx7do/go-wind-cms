import {defineConfig, presetIcons} from 'unocss'
import presetWeapp from 'unocss-preset-weapp'

export default defineConfig({
  presets: [
    presetIcons({
      scale: 1.2,
      warn: true,
      extraProperties: {
        'display': 'inline-block',
        'vertical-align': '-0.125em',
      },
    }),
    presetWeapp(),
  ],
})
