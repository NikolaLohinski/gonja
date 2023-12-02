{% import "macro.tpl" as m -%}
{% from "macro.tpl" import input, input as input2 -%}
<p>{{ m.input('username') }}</p>
<p>{{ m.input('username', simple.name) }}</p>
<p>{{ input('username', type='password') }}</p>
<p>{{ input2('username', type='password') }}</p>
