import React, { useEffect, useState } from 'react'

type Article = {
  id: number
  title: string
  summary: string
  author: string
  date: string
  url: string
  image?: string
}

type CacheRecord = {
  category: string
  ts: number
  articles: Article[]
}

const DB_NAME = 'newsapp-cache'
const DB_STORE = 'articles'
const CACHE_TTL_MS = 60 * 60 * 1000 // 1 hour

function openDb(): Promise<IDBDatabase> {
  return new Promise((resolve, reject) => {
    const req = indexedDB.open(DB_NAME, 1)
    req.onupgradeneeded = () => {
      const db = req.result
      if (!db.objectStoreNames.contains(DB_STORE)) {
        db.createObjectStore(DB_STORE, { keyPath: 'category' })
      }
    }
    req.onsuccess = () => resolve(req.result)
    req.onerror = () => reject(req.error)
  })
}

async function getCached(category: string): Promise<CacheRecord | null> {
  const db = await openDb()
  return new Promise((resolve, reject) => {
    const tx = db.transaction(DB_STORE, 'readonly')
    const store = tx.objectStore(DB_STORE)
    const r = store.get(category)
    r.onsuccess = () => resolve(r.result || null)
    r.onerror = () => reject(r.error)
  })
}

async function setCached(category: string, articles: Article[]) {
  const db = await openDb()
  return new Promise<void>((resolve, reject) => {
    const tx = db.transaction(DB_STORE, 'readwrite')
    const store = tx.objectStore(DB_STORE)
    const rec: CacheRecord = { category, ts: Date.now(), articles }
    const r = store.put(rec)
    r.onsuccess = () => resolve()
    r.onerror = () => reject(r.error)
  })
}

export default function App() {
  const [articles, setArticles] = useState<Article[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let mounted = true
    async function load() {
      setLoading(true)
      try {
        const cached = await getCached('general')
        if (cached && Date.now() - cached.ts < CACHE_TTL_MS) {
          if (mounted) {
            setArticles(cached.articles)
            setLoading(false)
          }
          // refresh in background
          fetchAndUpdate(false)
          return
        }
      } catch (e) {
        console.warn('idb read failed', e)
      }
      await fetchAndUpdate(true)
    }

    async function fetchAndUpdate(setBusy: boolean) {
      if (setBusy) setLoading(true)
      try {
        const res = await fetch('/api/articles')
        if (!res.ok) throw new Error(res.statusText)
        const data: Article[] = await res.json()
        if (mounted) setArticles(data || [])
        try {
          await setCached('general', data || [])
        } catch (e) {
          console.warn('idb write failed', e)
        }
      } catch (err) {
        console.error(err)
        if (mounted) setArticles([])
      } finally {
        if (mounted) setLoading(false)
      }
    }

    load()
    return () => {
      mounted = false
    }
  }, [])

  return (
    <div className="app">
      <header className="header">
        <div className="header-content">
          <div className="logo">
            <span className="logo-icon">📈</span>
            <h1>Stock Market News</h1>
          </div>
          <p className="subtitle">Real-time financial news from Finnhub</p>
        </div>
      </header>

      <div className="container">
        {loading ? (
          <div className="loading">
            <div className="spinner"></div>
            <p>Loading market news...</p>
          </div>
        ) : articles.length === 0 ? (
          <div className="empty">
            <p>No news available</p>
          </div>
        ) : (
          <main className="news-grid">
            {articles.map((a) => (
              <article key={a.id} className="card">
                {a.image && (
                  <div className="card-image">
                    <img src={a.image} alt={a.headline} />
                  </div>
                )}
                <div className="card-content">
                  <h2><a href={a.url} target="_blank" rel="noreferrer">{a.headline}</a></h2>
                  <p className="summary">{a.summary}</p>
                  <div className="meta">
                    <span className="source">{a.source}</span>
                    <span className="date">{new Date(a.datetime * 1000).toLocaleDateString()}</span>
                  </div>
                </div>
              </article>
            ))}
          </main>
        )}
      </div>

      <footer className="footer">
        <p>Market data provided by <strong>Finnhub</strong> • Built with React</p>
      </footer>
    </div>
  )
}
