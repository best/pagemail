'use client'

import { useState } from 'react'

interface ScrapeFormData {
  url: string
  email: string
  format: 'html' | 'pdf' | 'screenshot'
}

export function ScrapeForm() {
  const [formData, setFormData] = useState<ScrapeFormData>({
    url: '',
    email: '',
    format: 'html'
  })
  const [isLoading, setIsLoading] = useState(false)
  const [message, setMessage] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setMessage('')

    try {
      const response = await fetch('/api/v1/pages/scrape', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      })

      const result = await response.json()
      
      if (response.ok) {
        setMessage('请求已提交，文件将发送到您的邮箱！')
        setFormData({ url: '', email: '', format: 'html' })
      } else {
        setMessage(`错误：${result.error || '请求失败'}`)
      }
    } catch (error) {
      setMessage('网络错误，请稍后重试')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="bg-white shadow-lg rounded-lg p-6">
      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label htmlFor="url" className="block text-sm font-medium text-gray-700 mb-2">
            网页URL
          </label>
          <input
            type="url"
            id="url"
            required
            value={formData.url}
            onChange={(e) => setFormData({ ...formData, url: e.target.value })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="https://example.com"
          />
        </div>

        <div>
          <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
            接收邮箱
          </label>
          <input
            type="email"
            id="email"
            required
            value={formData.email}
            onChange={(e) => setFormData({ ...formData, email: e.target.value })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="your@email.com"
          />
        </div>

        <div>
          <label htmlFor="format" className="block text-sm font-medium text-gray-700 mb-2">
            输出格式
          </label>
          <select
            id="format"
            value={formData.format}
            onChange={(e) => setFormData({ ...formData, format: e.target.value as 'html' | 'pdf' | 'screenshot' })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="html">HTML 文件</option>
            <option value="pdf">PDF 文件</option>
            <option value="screenshot">网页截图</option>
          </select>
        </div>

        <button
          type="submit"
          disabled={isLoading}
          className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isLoading ? '处理中...' : '开始抓取'}
        </button>
      </form>

      {message && (
        <div className={`mt-4 p-3 rounded-md ${
          message.includes('错误') || message.includes('失败') 
            ? 'bg-red-100 text-red-700' 
            : 'bg-green-100 text-green-700'
        }`}>
          {message}
        </div>
      )}
    </div>
  )
}