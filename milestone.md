# Amazing Outfits — Project Milestones

A simple sneaker e-commerce website.

**Brand:** royal blue `#1E6FD9`, orange `#F28C1B`, cream `#F3EFE2`, dark green-black `#14211B`. Bold, sporty condensed type.

**Suggested stack:** Frontend (HTML/CSS/JS or React) · Backend (Go) · PostgreSQL · Cloudinary for images (or Supabase Storage / S3).

---

## Image upload strategy (the key decision)

- Only the **admin** uploads images, never customers.
- **Never store image files in the database.** Store them in object storage (Cloudinary is the easiest, with a free tier) and save only the **URL** in PostgreSQL.
- Flow: Admin form → backend validates (type: jpg/png/webp, max ~5 MB) → upload to Cloudinary → save returned URL in the `product_images` table.
- Cloudinary auto-resizes and converts to WebP, so thumbnails and product pages load fast.
- Each product gets 3–5 images (front, side, back, sole). The first one is the main image.
- Early on, skip uploads entirely: use placeholder image URLs while building the UI.

---

## Milestone 1 — Planning & Design
- [ ] Finalize brand colors, fonts, logo files
- [ ] Sketch pages: Home, Shop, Product, Cart, Checkout, Admin
- [ ] Define product fields: name, brand, price, sizes, colors, description, images, stock
- [ ] Choose image hosting (Cloudinary recommended)

## Milestone 2 — Frontend UI/UX (static, mock data)
- [ ] Header, footer, navigation, mobile menu
- [ ] Home page: hero, featured sneakers, categories
- [ ] Shop page: product grid, filters (brand, size, price), sort
- [ ] Product page: image gallery, size selector, add to cart
- [ ] Cart page and checkout form layout
- [ ] Fully responsive (mobile first) and accessible

## Milestone 3 — Backend & Database
- [ ] Set up server and PostgreSQL
- [ ] Tables: `products`, `product_images`, `sizes/stock`, `users`, `orders`, `order_items`
- [ ] REST API: list products, get product, create/update/delete (admin)
- [ ] Connect frontend to the API (replace mock data)

## Milestone 4 — Image Upload
- [ ] Create Cloudinary (or chosen storage) account and API keys (in `.env`, never in Git)
- [ ] Admin upload endpoint with file type and size validation
- [ ] Save image URLs in the database
- [ ] Multi-image upload with preview, reorder and delete in the admin form
- [ ] Lazy-load images on the frontend

## Milestone 5 — Cart, Checkout & Orders
- [ ] Cart (localStorage or server-side)
- [ ] Checkout and order creation
- [ ] Payment: M-Pesa (Daraja API) and/or card; start with "pay on delivery" if needed
- [ ] Order confirmation page and email

## Milestone 6 — Accounts & Admin Dashboard
- [ ] Customer sign up / login
- [ ] Order history
- [ ] Admin login, product management, order management

## Milestone 7 — Testing & Launch
- [ ] Test on phones and common browsers
- [ ] Performance check (image sizes, Lighthouse)
- [ ] Security check (input validation, auth, secrets)
- [ ] Deploy (frontend + backend + database) and connect a domain
- [ ] Collect feedback and plan v2 (wishlist, reviews, discounts)
