'use client'

interface StatsCardProps {
  label: string
  value: string | number
  subValue?: string
  progress?: number
  isBlue?: boolean
}

export function StatsCard({ label, value, subValue, progress, isBlue }: StatsCardProps) {
  return (
    <div className={`
      px-5 py-4 rounded-[2rem] border relative overflow-hidden w-full sm:w-[180px] transition-all hover:scale-105
      ${isBlue ? 'bg-blue-600 text-white border-blue-500 shadow-xl shadow-blue-500/10' : 'bg-white text-slate-900 border-slate-100 shadow-sm'}
    `}>
      <div className="relative z-10 flex flex-col items-center sm:items-start">
        <div className={`text-[9px] font-black uppercase tracking-[0.2em] mb-1 ${isBlue ? 'text-blue-100' : 'text-slate-400'}`}>
          {label}
        </div>
        <div className="flex items-baseline gap-1.5">
          <div className="text-3xl font-black italic">{value}</div>
          {subValue && <div className={`text-[9px] font-bold ${isBlue ? 'text-blue-100' : 'text-slate-400'}`}>{subValue}</div>}
        </div>
      </div>
      
      {progress !== undefined && (
        <div 
          className="absolute bottom-0 left-0 h-1 bg-white/20 transition-all duration-1000"
          style={{ width: `${progress}%` }}
        />
      )}
    </div>
  )
}
