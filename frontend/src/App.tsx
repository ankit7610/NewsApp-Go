import React, { useEffect, useState } from 'react'

type Article = {
  id: number
  title: string
  summary: string
  author: string
  date: string
  url: string
}

export default function App() {
  const [articles, setArticles] = useState<Article[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/articles')
      .then((res) => res.json())
      .then((data) => setArticles(data))
      .catch((err) => console.error(err))
      .finally(() => setLoading(false))
  }, [])

  return (
    <div className="container">
      <header>
        <h1>News</h1>
      </header>

      {loading ? (
        <p>Loading...</p>
      ) : (
        <main>
          {articles.map((a) => (
            <article key={a.id} className="card">
              <h2><a href={a.url} target="_blank" rel="noreferrer">{a.title}</a></h2>
              <p className="summary">{a.summary}</p>
              <p className="meta">{a.author} — {new Date(a.date).toLocaleDateString()}</p>
            </article>
          ))}
        </main>
      )}

      <footer>
        <small>Built with Go + TypeScript</small>
      </footer>
    </div>
  )
}
