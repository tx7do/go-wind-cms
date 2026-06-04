import {View, Text} from '@tarojs/components';
import {XIcon} from '@/plugins/xicon';
import {getCategoryName} from '@/api/hooks/category';
import {contentservicev1_Category} from '@/api/generated/app/service/v1';
import {useTranslations} from '@/lib/next-intl-compat';

const HomeCategoryCard: React.FC<{
    category: contentservicev1_Category;
    onClick?: (id: number) => void;
    mobileCompact?: boolean;
}> = ({category, onClick, mobileCompact = false}) => {
    const t = useTranslations('page.home');

    const handleClick = () => {
        if (!category?.id) return;
        onClick?.(category.id);
    };

    const getIconName = (icon?: string): string => {
        if (!icon) return 'carbon:folder';
        return icon.includes(':') ? icon : `carbon:${icon}`;
    };

    // ---- 横向滑动卡片（首页分类区） ----
    if (mobileCompact) {
        return (
            <View
              className='flex-shrink-0 w-[180rpx] flex flex-col items-center rounded bg-cardBg p-[20rpx]'
              onClick={handleClick}
              hoverClass='tap-active'
            >
                <View
                  className='flex items-center justify-center w-[72rpx] h-[72rpx] rounded mb-[12rpx]'
                  style={{backgroundColor: 'rgba(22,119,255,0.08)'}}
                >
                    <XIcon name={getIconName(category.icon)} size={28} className='text-primary' />
                </View>
                <Text
                  className='text-tips font-bold text-textMain text-center mb-[4rpx]'
                  style={{
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap',
                        maxWidth: '140rpx',
                    }}
                >
                    {getCategoryName(category)}
                </Text>
                <Text className='text-tips text-textThird'>
                    <Text className='text-primary'>{category.postCount || 0}</Text> {t('articles_unit')}
                </Text>
            </View>
        );
    }

    // ---- 默认列表卡片 ----
    return (
        <View
          className='flex flex-row items-center rounded bg-cardBg p-[24rpx] min-h-touch'
          onClick={handleClick}
          hoverClass='tap-active'
        >
            <View
              className='flex items-center justify-center w-[72rpx] h-[72rpx] rounded mr-[20rpx] flex-shrink-0'
              style={{backgroundColor: 'rgba(22,119,255,0.08)'}}
            >
                <XIcon name={getIconName(category.icon)} size={28} className='text-primary' />
            </View>
            <View className='flex-1 min-w-0'>
                <Text className='text-body font-bold text-textMain'>{getCategoryName(category)}</Text>
                <Text className='text-tips text-textThird'>
                    <Text className='text-primary'>{category.postCount || 0}</Text> {t('articles_unit')}
                </Text>
            </View>
        </View>
    );
};

export default HomeCategoryCard;
