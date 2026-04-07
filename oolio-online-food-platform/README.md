# Oolio Online Food Platform 🍰

A modern, responsive, and high-performance food ordering frontend built with React, TypeScript, and Vite. This platform offers a seamless dessert browsing and ordering experience with real-time cart updates and secure coupon validation.

## 🚀 Features

- **Dynamic Product Grid**: Browse a wide selection of desserts with optimized loading states and responsive layouts.
- **Real-time Cart**: fully interactive cart with instant quantity updates, removals, and total calculations.
- **Coupon Validation**: Integrated coupon system that validates promotional codes with the backend before finalizing orders.
- **Order Confirmation**: Beautiful modal-based order summary with persistent state management.
- **Responsive Design**: Mobile-first approach with a sticky sidebar cart for desktop users.
- **Secure Networking**: Centralized API handling with Axios interceptors for automatic API Key injection and CORS compliance.

## 🛠️ Tech Stack

- **Framework**: [React 19](https://react.dev/) + [Vite 8](https://vitejs.dev/)
- **Language**: [TypeScript](https://www.typescriptlang.org/)
- **State Management**: [Zustand](https://github.com/pmndrs/zustand) (Optimized with `useShallow`)
- **Data Fetching**: [TanStack Query v5](https://tanstack.com/query/latest) (React Query)
- **Styling**: [Tailwind CSS v4](https://tailwindcss.com/) + [Shadcn UI](https://ui.shadcn.com/)
- **Icons**: [Lucide React](https://lucide.dev/)
- **Routing**: [React Router v7](https://reactrouter.com/)

## 📦 Project Structure

```text
src/
├── components/
│   ├── pages/         # Page-specific components (Home, Product, Cart)
│   └── ui/            # Reusable Shadcn UI primitives
├── hooks/             # Custom React hooks (API integration)
├── lib/               # Utility libraries (Axios configuration)
├── pages/             # Main page views
├── services/          # API service definitions
├── store/             # Zustand global state (Cart & Checkout)
└── types/             # TypeScript interfaces and API models
```

## ⚙️ Getting Started

### Prerequisites

- Node.js (v18+)
- npm or yarn

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/azar-writes-code/oolio-assignment
   cd oolio-online-food-platform
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Configure Environment Variables:
   Create a `.env` file in the root directory:
   ```env
   VITE_API_URL=http://localhost:8080/api/v1
   VITE_API_KEY=apitest
   ```

### Development

Run the development server:
```bash
npm run dev
```

The app will be available at `http://localhost:5173`.

## 🔗 Backend Integration

This frontend is designed to work seamlessly with the **Oolio Product Service** (Go/Gin).
- **Default Port**: `8080`
- **Endpoints Used**: 
  - `GET /products/`: Fetch dessert list
  - `POST /order/`: Place a new order
  - `POST /order/validate-coupon`: Validate promotional codes

## 📜 License

Copyright (c) 2026 Azar. All rights reserved.

This project is **PROPRIETARY and CONFIDENTIAL**. Unauthorized use, reproduction, or distribution is strictly prohibited without express written permission from the copyright owner.
