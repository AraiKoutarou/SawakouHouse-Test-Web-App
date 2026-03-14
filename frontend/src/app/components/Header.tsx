'use client'

import { StatsCard } from './StatsCard'

interface HeaderProps {
  totalSpots: number
  completionRate: number
  prefecturesCount: number
}

export function Header({ totalSpots, completionRate, prefecturesCount }: HeaderProps) {
  return (
    <header className="bg-white border-b border-slate-100 py-8 md:py-12 mb-8 shadow-sm relative overflow-hidden">
      <div className="absolute top-0 right-0 w-64 h-64 bg-blue-50/50 rounded-full -translate-y-1/2 translate-x-1/2 blur-3xl pointer-events-none"></div>
      
      <div className="max-w-6xl mx-auto px-6">
        <div className="flex flex-col lg:flex-row items-center justify-between gap-8 md:gap-10">
          <div className="text-center lg:text-left">
            <div className="inline-block px-3 py-1 bg-blue-50 text-blue-600 rounded-full text-[9px] font-black uppercase tracking-[0.2em] mb-4 border border-blue-100">
              Personal Travel Explorer 🌎
            </div>
            <h1 className="text-4xl md:text-5xl lg:text-6xl font-black text-slate-900 tracking-tighter italic leading-tight">
              MEMORY <span className="text-blue-600">GLOW</span> MAP
            </h1>
            <p className="text-slate-400 font-bold mt-3 text-base md:text-lg max-w-xl mx-auto lg:mx-0 leading-relaxed">
              日本中を巡り、地図を思い出で染め上げよう。🌈
            </p>
          </div>
          
          <div className="flex flex-col sm:flex-row justify-center gap-4 w-full lg:w-auto">
            <StatsCard label="Total Spots" value={totalSpots} subValue="Places" />
            <StatsCard 
              label="Japan Completion" 
              value={`${completionRate}%`} 
              subValue={`${prefecturesCount} / 47 Prefs`} 
              progress={completionRate}
              isBlue={true} 
            />
          </div>
        </div>
      </div>
    </header>
  )
}
