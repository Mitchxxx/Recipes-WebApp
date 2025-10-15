# Recipes Web Frontend

React-based frontend application for the Recipes WebApp.

## Overview

This is a single-page application (SPA) built with React that provides a user interface for browsing and managing recipes. The application features Auth0 authentication integration and communicates with the backend API for data operations.

## Technology Stack

- **React 19.1**: UI library
- **Auth0 React SDK**: Authentication provider
- **Bootstrap 5.3**: CSS framework for styling
- **React Scripts**: Build tooling (Create React App)
- **Testing Library**: Component testing utilities

## Features

- Browse recipes from the backend API
- User authentication via Auth0
- Responsive design with Bootstrap
- Recipe cards with images and details
- Protected routes (authentication required)

## Prerequisites

- Node.js 14 or higher
- npm 6 or higher
- Running instance of the Recipes API (default: `http://localhost:8080`)

## Getting Started

This project was bootstrapped with [Create React App](https://github.com/facebook/create-react-app).

### Installation

1. **Install dependencies**
   ```bash
   npm install
   ```

2. **Configure Auth0**
   
   Update the Auth0 configuration in `src/index.js` with your Auth0 application credentials:
   ```javascript
   <Auth0Provider
     domain="your-auth0-domain.auth0.com"
     clientId="your-auth0-client-id"
     redirectUri={window.location.origin}
   >
   ```

3. **Configure API Endpoint**
   
   The application expects the API to be running at `http://localhost:8080`. If your API is hosted elsewhere, update the fetch URL in `src/App.js`:
   ```javascript
   fetch('http://your-api-url/recipes')
   ```

## Available Scripts

In the project directory, you can run:

### `npm start`

Runs the app in the development mode.\
Open [http://localhost:3000](http://localhost:3000) to view it in your browser.

**Note**: Make sure the backend API is running at `http://localhost:8080` before starting the frontend.

The page will reload when you make changes.\
You may also see any lint errors in the console.

### `npm test`

Launches the test runner in the interactive watch mode.\
See the section about [running tests](https://facebook.github.io/create-react-app/docs/running-tests) for more information.

### `npm run build`

Builds the app for production to the `build` folder.\
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.\
Your app is ready to be deployed!

See the section about [deployment](https://facebook.github.io/create-react-app/docs/deployment) for more information.

### `npm run eject`

**Note: this is a one-way operation. Once you `eject`, you can't go back!**

If you aren't satisfied with the build tool and configuration choices, you can `eject` at any time. This command will remove the single build dependency from your project.

Instead, it will copy all the configuration files and the transitive dependencies (webpack, Babel, ESLint, etc) right into your project so you have full control over them. All of the commands except `eject` will still work, but they will point to the copied scripts so you can tweak them. At this point you're on your own.

You don't have to ever use `eject`. The curated feature set is suitable for small and middle deployments, and you shouldn't feel obligated to use this feature. However we understand that this tool wouldn't be useful if you couldn't customize it when you are ready for it.

## Project Structure

```
recipes-web/
├── public/              # Static files
│   ├── index.html      # HTML template
│   └── favicon.ico     # App icon
├── src/
│   ├── App.js          # Main application component
│   ├── App.css         # App styles
│   ├── Recipe.js       # Recipe card component
│   ├── Recipe.css      # Recipe card styles
│   ├── Navbar.js       # Navigation bar with Auth0 integration
│   ├── Profile.js      # User profile component
│   ├── index.js        # App entry point with Auth0 provider
│   └── index.css       # Global styles
├── Dockerfile          # Container image definition
├── docker-compose.yml  # Docker orchestration
└── package.json        # Dependencies and scripts
```

## Components

### App
The main application component that fetches recipes from the API and renders them in a grid layout.

### Navbar
Navigation bar component with Auth0 authentication integration. Shows login button or user profile based on authentication status.

### Recipe
Card component that displays individual recipe information including title, image, and details.

### Profile
User profile component displayed when authenticated, showing user information from Auth0.

## Docker Deployment

### Build the Docker image
```bash
docker build -t recipes-web .
```

### Run with Docker Compose
```bash
docker-compose up -d
```

## Environment Configuration

### Production Build
When building for production, you may need to update the API URL. Modify `src/App.js`:

```javascript
const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';
fetch(`${API_URL}/recipes`)
```

Then set the environment variable:
```bash
REACT_APP_API_URL=https://your-api-domain.com npm run build
```

## Authentication

The app uses Auth0 for authentication. Users can:
- Log in with Auth0 credentials
- View their profile information
- Access protected features (when implemented)

Auth0 configuration is in `src/index.js`.

## Learn More

You can learn more in the [Create React App documentation](https://facebook.github.io/create-react-app/docs/getting-started).

To learn React, check out the [React documentation](https://reactjs.org/).

### Code Splitting

This section has moved here: [https://facebook.github.io/create-react-app/docs/code-splitting](https://facebook.github.io/create-react-app/docs/code-splitting)

### Analyzing the Bundle Size

This section has moved here: [https://facebook.github.io/create-react-app/docs/analyzing-the-bundle-size](https://facebook.github.io/create-react-app/docs/analyzing-the-bundle-size)

### Making a Progressive Web App

This section has moved here: [https://facebook.github.io/create-react-app/docs/making-a-progressive-web-app](https://facebook.github.io/create-react-app/docs/making-a-progressive-web-app)

### Advanced Configuration

This section has moved here: [https://facebook.github.io/create-react-app/docs/advanced-configuration](https://facebook.github.io/create-react-app/docs/advanced-configuration)

### Deployment

This section has moved here: [https://facebook.github.io/create-react-app/docs/deployment](https://facebook.github.io/create-react-app/docs/deployment)

### `npm run build` fails to minify

This section has moved here: [https://facebook.github.io/create-react-app/docs/troubleshooting#npm-run-build-fails-to-minify](https://facebook.github.io/create-react-app/docs/troubleshooting#npm-run-build-fails-to-minify)
