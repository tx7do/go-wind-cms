import {View, Text} from '@tarojs/components';
import {XIcon} from '@/plugins/xicon';
import {getCategoryName, getCategoryDescription} from '@/api/hooks/category';
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

    const name = getCategoryName(category);
    const desc = getCategoryDescription?.(category);
    const count = category.postCount || 0;

    // ---- 横向滑动大卡片（首页分类区） ----
    if (mobileCompact) {
        return (
            <View
              className='flex-shrink-0 rounded-[16rpx] bg-cardBg overflow-hidden'
              style={{width: '260rpx'}}
              onClick={handleClick}
              hoverClass='tap-active'
            >
                {/* 顶部图标区 - 增强背景色 */}
                <View
                  className='flex items-center justify-center'
                  style={{
                        height: '120rpx',
                        backgroundColor: 'rgba(22,119,255,0.1)',
                    }}
                >
                    <XIcon name={getIconName(category.icon)} size={40} className='text-primary' />
                </View>
                {/* 文字区 */}
                <View className='p-[20rpx]'>
                    <Text
                      className='text-desc font-bold text-textMain mb-[8rpx]'
                      numberOfLines={1}
                      style={{
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            display: '-webkit-box',
                            WebkitLineClamp: 1,
                            WebkitBoxOrient: 'vertical',
                        }}
                    >
                        {name}
                    </Text>
                    <Text className='text-tips text-textThird'>
                        <Text className='text-primary font-medium'>{count}</Text>
                        <Text className='text-tips text-textThird'>篇</Text>
                    </Text>
                </View>
            </View>
        );
    }

    // ---- 默认列表卡片（横排） ----
    return (
        <View
          className='flex flex-row items-center rounded bg-cardBg p-[24rpx]'
          style={{minHeight: '88rpx'}}
          onClick={handleClick}
          hoverClass='tap-active'
        >
            <View
              className='flex items-center justify-center rounded flex-shrink-0'
              style={{
                    width: '72rpx',
                    height: '72rpx',
                    marginRight: '20rpx',
                    backgroundColor: 'rgba(22,119,255,0.08)',
                }}
            >
                <XIcon name={getIconName(category.icon)} size={28} className='text-primary' />
            </View>
            <View className='flex-1 min-w-0'>
                <Text className='text-body font-bold text-textMain'>{name}</Text>
                {desc && (
                    <Text
                      className='text-tips text-textSec'
                      style={{
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            whiteSpace: 'nowrap',
                        }}
                    >
                        {desc}
                    </Text>
                )}
            </View>
            <Text className='text-tips text-textThird flex-shrink-0' style={{marginLeft: '16rpx'}}>
                <Text className='text-primary'>{count}</Text> {t('articles_unit')}
            </Text>
        </View>
    );
};

export default HomeCategoryCard;
