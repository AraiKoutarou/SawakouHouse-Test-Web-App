'use client'

import { CategoryBadge } from './ui/CategoryBadge'

interface Memory {
  id: number
  title: string
  address: string
  prefecture: string
  category: string
  comment: string
  color: string
}

export function MemoryCard({ memory }: { memory: Memory }) {
  return (
    <div className="pop-card bg-white p-5 md:p-7 rounded-[2rem] md:rounded-[2.5rem] shadow-sm flex flex-col h-full border border-slate-100 group">
      <div className="flex items-start justify-between mb-2 overflow-hidden">
        <div className="flex items-center gap-2 md:gap-3 flex-grow min-w-0 pr-2">
          <div 
            className="w-4 h-4 md:w-5 md:h-5 rounded-full border-2 border-white shadow-md flex-shrink-0" 
            style={{ 
              backgroundColor: memory.color, 
              boxShadow: `0 0 10px ${memory.color}` 
            }} 
          />
          <h3 className="font-black text-lg md:text-xl text-slate-800 truncate group-hover:text-blue-600 transition-colors italic">
            {memory.title}
          </h3>
        </div>
        <div className="flex-shrink-0">
          <CategoryBadge categoryId={memory.category} />
        </div>
      </div>

      <p className="text-[9px] md:text-[10px] text-slate-400 font-bold uppercase mb-4 md:mb-5 truncate pl-6 md:pl-8">
        {memory.address}
      </p>

      <div className="bg-slate-50/50 p-4 md:p-5 rounded-2xl md:rounded-3xl mb-4 md:mb-6 flex-grow border border-slate-50 group-hover:bg-white group-hover:border-slate-100 transition-all duration-300">
        <p className="text-slate-600 text-[14px] md:text-[15px] leading-relaxed italic font-medium line-clamp-4">
          “{memory.comment || "素敵な時間の欠片がここに...✨"}”
        </p>
      </div>

      <div className="flex items-center justify-between text-[9px] md:text-[10px] font-black uppercase tracking-widest pt-4 md:pt-5 border-t border-slate-50 mt-auto">
        <div className="flex items-center gap-2">
          <span className="text-blue-600 bg-blue-50 px-2 py-0.5 rounded-lg">{memory.prefecture}</span>
          <span className="text-slate-200">•</span>
          <span className="text-slate-400 font-medium">{new Date(memory.id).toLocaleDateString()}</span>
        </div>
      </div>
    </div>
  )
}
