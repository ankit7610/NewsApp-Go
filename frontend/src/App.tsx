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

type Category = 'general' | 'forex' | 'crypto' | 'merger'

export default function App() {
  const [articles, setArticles] = useState<Article[]>([])
  const [loading, setLoading] = useState(true)
  const [category, setCategory] = useState<Category>('general')

  useEffect(() => {
    setLoading(true)
    const token = (import.meta as any).env.VITE_FINNHUB_API_KEY || ''
    if (!token) {
      console.warn('VITE_FINNHUB_API_KEY not set; Finnhub requests may fail or be rate-limited.')
    }
    const url = `https://finnhub.io/api/v1/news?category=${category}&token=${token}`
    fetch(url)
      .then((res) => {
        if (!res.ok) throw new Error(res.statusText)
        return res.json()
      })
      .then((data: any[]) => {
        const mapped = (data || []).map((fa) => ({
          id: fa.id,
          title: fa.headline,
          summary: fa.summary,
          author: fa.source,
          date: new Date((fa.datetime || 0) * 1000).toISOString(),
          url: fa.url,
          image: fa.image,
        }))
        setArticles(mapped)
      })
      .catch((err) => {
        console.error(err)
        setArticles([])
      })
      .finally(() => setLoading(false))
  }, [category])

  const categories: { value: Category; label: string; icon: string }[] = [
    { value: 'general', label: 'Market News', icon: '📰' },
    { value: 'forex', label: 'Forex', icon: '💱' },
    { value: 'crypto', label: 'Crypto', icon: '₿' },
    { value: 'merger', label: 'M&A', icon: '🤝' },
  ]

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
        <div className="filters">
          {categories.map((cat) => (
            <button
              key={cat.value}
              className={`filter-btn ${category === cat.value ? 'active' : ''}`}
              onClick={() => setCategory(cat.value)}
            >
              <span className="icon">{cat.icon}</span>
              {cat.label}
            </button>
          ))}
        </div>

        {loading ? (
          <div className="loading">
            <div className="spinner"></div>
            <p>Loading market news...</p>
          </div>
        ) : articles.length === 0 ? (
          <div className="empty">
            <p>No news available for this category</p>
          </div>
        ) : (
          <main className="news-grid">
            {articles.map((a) => (
              <article key={a.id} className="card">
                {a.image && (
                  <div className="card-image">
                    <img src={a.image} alt={a.title} />
                  </div>
                )}
                <div className="card-content">
                  <h2><a href={a.url} target="_blank" rel="noreferrer">{a.title}</a></h2>
                  <p className="summary">{a.summary}</p>
                  <div className="meta">
                    <span className="source">{a.author}</span>
                    <span className="date">{new Date(a.date).toLocaleDateString()}</span>
                  </div>
                </div>
              </article>
            ))}
          </main>
        )}
      </div>

      <footer className="footer">
        <p>Market data provided by <strong>Finnhub</strong> • Built with Go + React</p>
      </footer>
    </div>
  )
}
