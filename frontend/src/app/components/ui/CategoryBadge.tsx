'use client'

import { CATEGORIES, CategoryInfo } from '@/app/lib/categories'

export function CategoryBadge({ categoryId }: { categoryId: string }) {
  const info = CATEGORIES[categoryId] || CATEGORIES.other;
  
  return (
    <div 
      className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-widest border transition-all"
      style={{ 
        backgroundColor: info.bg, 
        color: info.color, 
        borderColor: info.border 
      }}
    >
      <span>{info.emoji}</span>
      <span>{info.label}</span>
    </div>
  )
}
