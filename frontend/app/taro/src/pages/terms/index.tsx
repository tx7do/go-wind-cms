import {useTranslation} from 'react-i18next';
import {View, Text} from '@tarojs/components';

import '../../globals.css'; // 导入全局 CSS，确保 CSS 变量可用
import styles from './terms.scss';

export default function TermsPage() {
  const {t} = useTranslation('page.legal.terms');

  return (
    <View className={styles.infoPage}>
      {/* Hero Section */}
      <View className={styles.hero}>
        <View className={styles.heroContent}>
          <Text className={styles.heroTitle}>{t('title')}</Text>
          <Text className={styles.heroSubtitle}>{t('description')}</Text>
        </View>
      </View>

      {/* Content Section */}
      <View className={styles.infoCard}>
        <View className={styles.list}>
          <Text>{t('item_1')}</Text>
          <Text>{t('item_2')}</Text>
          <Text>{t('item_3')}</Text>
        </View>
      </View>
    </View>
  );
}
