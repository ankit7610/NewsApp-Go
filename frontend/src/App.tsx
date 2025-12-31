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

export default function App() {
  const [articles, setArticles] = useState<Article[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/articles')
      .then((res) => {
        if (!res.ok) throw new Error(res.statusText)
        return res.json()
      })
      .then((data: Article[]) => {
        setArticles(data || [])
      })
      .catch((err) => {
        console.error(err)
        setArticles([])
      })
      .finally(() => setLoading(false))
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
