'use client'

import { APIProvider } from '@vis.gl/react-google-maps'
import { useState, useEffect, useMemo } from 'react'
import GoogleMap from './components/Map'
import { Header } from './components/Header'
import { MemoryCard } from './components/MemoryCard'

interface Location {
  id: number
  place_id: string
  title: string
  address: string
  prefecture: string
  category: string
  comment: string
  color: string
  latitude: number
  longitude: number
}

const GOOGLE_MAPS_API_KEY = process.env.NEXT_PUBLIC_GOOGLE_MAPS_API_KEY || '';
const MAP_ID = process.env.NEXT_PUBLIC_MAP_ID || 'DEMO_MAP_ID';

export default function Home() {
  const [locations, setLocations] = useState<Location[]>([])

  const fetchLocations = async () => {
    try {
      const res = await fetch('http://localhost:8080/locations')
      if (res.ok) {
        const data = await res.json()
        setLocations(data || [])
      }
    } catch (err) {
      console.error('Fetch error:', err)
    }
  }

  useEffect(() => {
    fetchLocations()
  }, [])

  const addLocation = async (data: Partial<Location>) => {
    try {
      const res = await fetch('http://localhost:8080/locations', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
      })

      if (res.ok) {
        fetchLocations()
      }
    } catch (err) {
      console.error('Save error:', err)
    }
  }

  const visitedPrefectures = useMemo(() => {
    const prefs = new Set(locations.map(l => l.prefecture).filter(p => p !== ''))
    return Array.from(prefs)
  }, [locations])

  const completionRate = Math.round((visitedPrefectures.length / 47) * 100)

  return (
    <main className="min-h-screen pb-20 bg-slate-50/30">
      <Header 
        totalSpots={locations.length} 
        completionRate={completionRate} 
        prefecturesCount={visitedPrefectures.length} 
      />

      <div className="max-w-6xl mx-auto px-4 md:px-6">
        <section className="mb-16 md:mb-24">
          <APIProvider apiKey={GOOGLE_MAPS_API_KEY} libraries={['places']}>
            <GoogleMap locations={locations} mapId={MAP_ID} onAddLocation={addLocation} />
          </APIProvider>
        </section>

        <section>
          <div className="flex flex-col md:flex-row md:items-end justify-between gap-6 mb-8 md:mb-12">
            <div>
              <h2 className="text-3xl md:text-4xl font-black text-slate-900 italic tracking-tighter">JOURNAL</h2>
              <p className="text-slate-400 font-bold uppercase text-[9px] tracking-[0.3em] mt-1 md:mt-2">Your Magic Moments Timeline</p>
            </div>
            
            {visitedPrefectures.length > 0 && (
              <div className="flex flex-wrap gap-1.5 md:gap-2 justify-start md:justify-end">
                {visitedPrefectures.map(pref => (
                  <span key={pref} className="px-3 py-1 bg-white text-blue-600 text-[9px] font-black rounded-full border border-slate-100 shadow-sm transition-all hover:scale-110">
                    {pref}
                  </span>
                ))}
              </div>
            )}
          </div>
          
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 md:gap-10">
            {locations.length === 0 ? (
              <div className="col-span-full py-20 md:py-32 bg-white rounded-[2.5rem] md:rounded-[3.5rem] shadow-inner border-2 border-dashed border-slate-100 flex flex-col items-center justify-center text-center px-6">
                <span className="text-5xl md:text-7xl mb-6 animate-bounce">🛸</span>
                <p className="text-slate-400 font-bold text-lg md:text-2xl leading-relaxed">
                  まだ記録された思い出はありません。<br/>
                  <span className="text-blue-600">最初のお店を検索しましょう！</span>
                </p>
              </div>
            ) : (
              locations.map((loc) => (
                <MemoryCard key={loc.id} memory={loc} />
              ))
            )}
          </div>
        </section>
      </div>

      <footer className="mt-24 md:mt-40 text-center px-6">
        <div className="inline-flex items-center gap-4 px-6 md:px-8 py-3 md:py-4 bg-white rounded-full shadow-sm border border-slate-100">
          <span className="text-slate-300 font-black italic text-[9px] md:text-xs tracking-[0.2em]">&copy; 2026</span>
          <div className="w-1.5 h-1.5 bg-blue-500 rounded-full animate-ping"></div>
          <span className="text-slate-400 font-black italic text-[9px] md:text-xs tracking-[0.2em]">MEMORY GLOW EXPLORER</span>
        </div>
      </footer>
    </main>
  )
}
