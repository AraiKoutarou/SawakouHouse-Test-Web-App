'use client'

import { 
  Map, 
  AdvancedMarker, 
  InfoWindow, 
  MapMouseEvent,
  useMapsLibrary,
  useMap
} from '@vis.gl/react-google-maps'
import { useState, useCallback, useEffect, useRef } from 'react'
import { getCategoryFromType, CATEGORIES } from '@/app/lib/categories'

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

interface MapProps {
  locations: Location[]
  mapId?: string
  onAddLocation: (data: Partial<Location>) => void
}

// --- 入力フォームを独立したコンポーネントに分離 ---
// これにより、入力中のステート更新が地図全体に波及せず、IMEの変換が途切れなくなります。
interface MemoryFormProps {
  place: Partial<Location>;
  onSave: (data: { comment: string; color: string }) => void;
}

function MemoryForm({ place, onSave }: MemoryFormProps) {
  const [comment, setComment] = useState('')
  const [color, setColor] = useState('#3b82f6')
  const [isComposing, setIsComposing] = useState(false)

  // 初期色のセット（カテゴリーに応じた色）
  useEffect(() => {
    const cat = getCategoryFromType([place.category || ''])
    setColor(cat.color)
  }, [place.category])

  const handleSave = () => {
    if (isComposing) return // 変換中は何もしない
    onSave({ comment, color })
  }

  return (
    <div className="p-3 w-[260px] md:w-[300px] bg-white rounded-2xl md:rounded-[2rem]">
      <div className="flex items-center gap-2 mb-1.5">
        <span className="text-xl md:text-2xl">{CATEGORIES[place.category!]?.emoji || '📍'}</span>
        <h3 className="font-black text-lg md:text-xl text-blue-600 truncate flex-grow">{place.title}</h3>
      </div>
      <p className="text-[8px] text-slate-400 mb-3 truncate pl-8">{place.address || '緯度・経度から指定'}</p>
      
      <div className="space-y-3">
        <textarea
          placeholder="どんな思い出ですか？ ✍️"
          className="w-full p-3 bg-slate-50 border border-slate-100 focus:border-blue-300 rounded-xl md:rounded-2xl outline-none font-medium text-slate-700 h-24 resize-none transition-all text-xs"
          value={comment}
          onChange={(e) => setComment(e.target.value)}
          onCompositionStart={() => setIsComposing(true)}
          onCompositionEnd={() => setIsComposing(false)}
          onKeyDown={(e) => {
            if (e.key === 'Enter' && e.ctrlKey) handleSave() // Ctrl+Enter で保存
          }}
          autoFocus
        />

        <div className="flex items-center justify-between px-1">
          <span className="text-[9px] font-black text-slate-400 uppercase tracking-widest">Glow Color</span>
          <input
            type="color"
            className="w-10 h-8 rounded-lg cursor-pointer border-none shadow-sm"
            value={color}
            onChange={(e) => setColor(e.target.value)}
          />
        </div>

        <button
          onClick={handleSave}
          className="w-full py-3 md:py-4 bg-blue-600 text-white font-black rounded-xl md:rounded-[1.5rem] shadow-lg hover:bg-blue-700 active:scale-95 transition-all uppercase tracking-tighter text-sm"
        >
          Record Magic! 🪄
        </button>
      </div>
    </div>
  )
}

export default function GoogleMap({ locations, mapId, onAddLocation }: MapProps) {
  const map = useMap()
  const places = useMapsLibrary('places')
  
  const [newPlace, setNewPlace] = useState<Partial<Location> | null>(null)
  const [selectedLocation, setSelectedLocation] = useState<Location | null>(null)
  
  const searchInputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    if (!places || !map || !searchInputRef.current) return

    const autocomplete = new places.Autocomplete(searchInputRef.current, {
      fields: ['place_id', 'name', 'formatted_address', 'address_components', 'geometry', 'types'],
      componentRestrictions: { country: 'jp' }
    })

    autocomplete.addListener('place_changed', () => {
      const place = autocomplete.getPlace()
      if (!place.geometry || !place.geometry.location) return

      map.panTo(place.geometry.location)
      map.setZoom(17)

      let pref = ''
      if (place.address_components) {
        const prefComp = place.address_components.find(c => c.types.includes('administrative_area_level_1'))
        if (prefComp) pref = prefComp.long_name
      }
      
      const cat = getCategoryFromType(place.types)

      setNewPlace({
        place_id: place.place_id,
        title: place.name,
        address: place.formatted_address,
        prefecture: pref,
        category: cat.id,
        latitude: place.geometry.location.lat(),
        longitude: place.geometry.location.lng()
      })
    })
  }, [places, map])

  const handleMapClick = useCallback((e: MapMouseEvent) => {
    if (e.detail.latLng) {
      setNewPlace({
        title: '選択した地点',
        category: 'other',
        latitude: e.detail.latLng.lat,
        longitude: e.detail.latLng.lng,
      })
      setSelectedLocation(null)
    }
  }, [])

  const handleSaveMemory = (formData: { comment: string; color: string }) => {
    if (newPlace && newPlace.title) {
      onAddLocation({ ...newPlace, ...formData })
      setNewPlace(null)
    }
  }

  return (
    <div className="h-[500px] md:h-[650px] w-full rounded-[2.5rem] md:rounded-[3rem] shadow-2xl border-[6px] md:border-[12px] border-white overflow-hidden relative group bg-slate-50">
      <div className="absolute top-4 md:top-8 left-1/2 -translate-x-1/2 z-10 w-[95%] max-w-md">
        <div className="relative">
          <div className="absolute inset-y-0 left-4 flex items-center pointer-events-none text-lg">🔍</div>
          <input
            ref={searchInputRef}
            type="text"
            placeholder="思い出の場所を探してみよう...✨"
            className="w-full pl-12 pr-5 py-3.5 md:py-5 bg-white/95 backdrop-blur-md border-none rounded-2xl md:rounded-[2rem] shadow-xl text-slate-800 font-bold focus:ring-4 focus:ring-blue-400/30 outline-none transition-all placeholder-slate-400 text-sm"
          />
        </div>
      </div>

      <Map
        defaultCenter={{lat: 35.6895, lng: 139.6917}}
        defaultZoom={13}
        mapId={mapId || "DEMO_MAP_ID"}
        onClick={handleMapClick}
        disableDefaultUI={true}
      >
        {locations.map((loc) => (
          <AdvancedMarker
            key={loc.id}
            position={{lat: loc.latitude, lng: loc.longitude}}
            onClick={() => setSelectedLocation(loc)}
          >
            <div className="memory-glow" style={{ backgroundColor: loc.color, boxShadow: `0 0 30px 10px ${loc.color}`, color: loc.color }} />
            <div className="memory-core" style={{ color: loc.color }} />
          </AdvancedMarker>
        ))}

        {selectedLocation && (
          <InfoWindow
            position={{lat: selectedLocation.latitude, lng: selectedLocation.longitude}}
            onCloseClick={() => setSelectedLocation(null)}
          >
            <div className="p-2 text-slate-800 w-[200px] md:w-[240px]">
              <div className="flex items-center gap-2 mb-2">
                <span className="text-lg">{CATEGORIES[selectedLocation.category]?.emoji || '📍'}</span>
                <h3 className="font-bold text-base md:text-lg leading-tight italic truncate flex-grow" style={{color: selectedLocation.color}}>
                  {selectedLocation.title}
                </h3>
              </div>
              <p className="text-[9px] text-slate-400 mb-2 truncate leading-relaxed">{selectedLocation.address}</p>
              <div className="bg-slate-50 p-2.5 rounded-xl italic text-[12px] text-slate-600 font-medium leading-relaxed border border-slate-100 max-h-24 overflow-y-auto">
                “{selectedLocation.comment}”
              </div>
              <div className="mt-2.5 flex justify-between items-center">
                <div className="text-[8px] font-black text-blue-500 uppercase tracking-widest px-1.5 py-0.5 bg-blue-50 rounded">
                  {selectedLocation.prefecture}
                </div>
                <div className="text-[8px] text-slate-300 font-bold">{new Date(selectedLocation.id).toLocaleDateString()}</div>
              </div>
            </div>
          </InfoWindow>
        )}

        {newPlace && (
          <InfoWindow
            position={{lat: newPlace.latitude!, lng: newPlace.longitude!}}
            onCloseClick={() => setNewPlace(null)}
          >
            <MemoryForm 
              place={newPlace} 
              onSave={handleSaveMemory} 
            />
          </InfoWindow>
        )}
      </Map>
    </div>
  )
}
