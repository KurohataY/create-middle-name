$(function () {
  // ヒーローイメージの文字の要素を取得
  var hero = $("h1");
  // 取得した要素から文字だけを取得する
  var heroText = hero.text();
  // 取得したテキストを一文字づつ分割した配列に変換
  var heroSplit = heroText.split('');
  // 配列から文字ひとつひとつ span タグで囲む
  var heroSpan = "";
  for (var i = 0;i < heroSplit.length;i++) {
    heroSpan += (`<span>${heroSplit[i]}</span>`);
  }
  // span タグで囲った文字とspanタグを h1 に入れ直す
  $("h1").html(heroSpan);
  // span を透明（不透明度 0%）に
  // $("span").css({opacity:0});
  // 順番に span の不透明度を上げていく
  var n=0;
  function func(){
    // $("span").eq(n).css({opacity:'0.0'}).animate({opacity: '1'}, 1500);
    $("h1 span").eq(n).addClass('hero-anime');
    n++;
    var tid = setTimeout(function(){
      func();
    },100);
    if (n>heroSpan.length){
      clearTimeout(tid);
    }
  }
  func();

});
