import { SearchBar } from './components/SearchBar';
import './styles/global.css';

// Initialize the search bar when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
  new SearchBar();
  console.log('Memex search ready.');
});
