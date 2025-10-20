import http from 'k6/http';
import { sleep } from 'k6';
export const options = { vus: 20, duration: '1m' };

export default function () {
  http.get('http://host.docker.internal:8080/api/catalog/list');
  sleep(1);
}
