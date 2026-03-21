import {View} from '@tarojs/components'

import './index.scss'

interface Props {
  current: number
  total: number
  pageSize: number
  onChange: (page: number) => void
  showSizeChanger?: boolean
}

export default function Pagination({
                                     current,
                                     total,
                                     pageSize,
                                     onChange,
                                     showSizeChanger = false
                                   }: Props) {
  const totalPages = Math.ceil(total / pageSize)
  const pages: number[] = []

  // 生成页码数组
  for (let i = 1; i <= totalPages; i++) {
    pages.push(i)
  }

  const handlePrev = () => {
    if (current > 1) {
      onChange(current - 1)
    }
  }

  const handleNext = () => {
    if (current < totalPages) {
      onChange(current + 1)
    }
  }

  const handlePage = (page: number) => {
    onChange(page)
  }

  return (
    <View className='pagination'>
      <View
        className={`pagination-btn ${current === 1 ? 'disabled' : ''}`}
        onClick={handlePrev}
      >
        上一页
      </View>

      {pages.map((page) => (
        <View
          key={page}
          className={`pagination-item ${current === page ? 'active' : ''}`}
          onClick={() => handlePage(page)}
        >
          {page}
        </View>
      ))}

      <View
        className={`pagination-btn ${current === totalPages ? 'disabled' : ''}`}
        onClick={handleNext}
      >
        下一页
      </View>

      {showSizeChanger && (
        <View className='pagination-total'>
          共 {total} 条
        </View>
      )}
    </View>
  )
}
