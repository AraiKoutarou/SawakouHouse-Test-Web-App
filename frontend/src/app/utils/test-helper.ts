// テスト用：わざと改善の余地があるコードを含めています

export function calc(a: any, b: any) {
  return a + b;
}

export function getData(url: string) {
  fetch(url)
    .then(res => res.json())
    .then(data => {
      console.log(data);
      return data;
    });
}

export const processUser = (user: any) => {
  if(user) {
    if(user.name) {
      if(user.age > 0) {
        return user.name + " is " + user.age + " years old";
      }
    }
  }
  return "Unknown";
}

export function divide(x: number, y: number) {
  return x / y;
}

let globalCounter = 0;
export function incrementCounter() {
  globalCounter = globalCounter + 1;
  return globalCounter;
}
