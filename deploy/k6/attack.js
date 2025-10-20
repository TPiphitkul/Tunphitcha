import http from 'k6/http';
export const options = { vus: 100, duration: '1m' };

export default function () {
  http.post('http://host.docker.internal:8080/api/user/login', JSON.stringify({u:'a',p:'bad'}), {
    headers: { 'Content-Type':'application/json' }
  });
}
