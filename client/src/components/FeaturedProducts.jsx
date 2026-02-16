import React from 'react'
import ProductCard from './ProductCard'
import './FeaturedProducts.css'

function FeaturedProducts({ products }) {
  if (products.length === 0) return null

  return (
    <div className="featured-products">
      <div className="featured-grid">
        {products.map((product) => (
          <ProductCard key={product.id} product={product} />
        ))}
      </div>
    </div>
  )
}

export default FeaturedProducts

