// Harf: original Windows input helper. Build with .NET Framework C# compiler.
using System;
using System.IO;
using System.Text;
using System.Text.RegularExpressions;
using System.Diagnostics;
using System.Drawing;
using System.Runtime.InteropServices;
using System.Threading;
using System.Windows.Forms;
using Microsoft.Win32;

static class Program {
    [STAThread] static void Main(string[] args) {
        Application.EnableVisualStyles();
        Application.SetCompatibleTextRenderingDefault(false);
        if (Array.IndexOf(args,"--self-test")>=0) { Converter.Test(); return; }
        if (Path.GetFileNameWithoutExtension(Application.ExecutablePath).EndsWith("-Setup",StringComparison.OrdinalIgnoreCase)) { Install(); return; }
        bool first; using(var mutex=new Mutex(true,"Local\\Harf.InputHelper.v1",out first)) {
            if(!first){MessageBox.Show("Harf is already running. Open its tray icon.","Harf");return;}
            Application.Run(new HarfWindow());
        }
    }
    static void Install() {
        if(MessageBox.Show("Install Harf for this Windows user?\n\nThis installs an input helper, not a Windows keyboard layout.\nAutomatic replacements are OFF until you enable them.","Harf Setup",MessageBoxButtons.OKCancel,MessageBoxIcon.Information)!=DialogResult.OK)return;
        string folder=Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.LocalApplicationData),"Harf");
        try {
            Directory.CreateDirectory(folder);
            string target=Path.Combine(folder,"Harf.exe");
            File.Copy(Application.ExecutablePath,target,true);
            string start=Path.Combine(Environment.GetFolderPath(Environment.SpecialFolder.StartMenu),"Programs","Harf.url");
            File.WriteAllText(start,"[InternetShortcut]\r\nURL=file:///"+target.Replace('\\','/')+"\r\nIconFile="+target+"\r\nIconIndex=0\r\n");
            string uninstall=Path.Combine(folder,"Uninstall.cmd");
            File.WriteAllText(uninstall,"@echo off\r\necho Close Harf from its tray menu first.\r\npause\r\nreg delete HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run /v Harf /f >nul 2>&1\r\ndel \""+start+"\" >nul 2>&1\r\ndel \"%~dp0Harf.exe\"\r\nif exist \"%~dp0Harf.exe\" (echo Harf is still running. Close it and retry. & pause & exit /b 1)\r\ncd /d \"%TEMP%\"\r\nrmdir /s /q \""+folder+"\"\r\n",Encoding.Default);
            Process.Start(target);
        }catch(Exception ex){MessageBox.Show(ex.Message,"Installation failed",MessageBoxButtons.OK,MessageBoxIcon.Error);}
    }
}

static class Converter {
    const string Apostrophes="'`´ʻʼ‘’";
    static bool Letter(string s,int n){return n>=0&&n<s.Length&&char.IsLetter(s[n]);}
    public static string Convert(string input,bool reverse){
        return Regex.Replace(input,@"\S+",m=>{
            string s=m.Value;if(s.Contains("://")||s.Contains("@")||s.StartsWith("www.",StringComparison.OrdinalIgnoreCase))return s;
            if(reverse){bool caps=Regex.IsMatch(s,@"\p{Lu}")&&!Regex.IsMatch(s,@"\p{Ll}");return Regex.Replace(s,"[öÖğĞşŞçÇ]",x=>{
                switch(x.Value){case "ö":return "o‘";case "Ö":return "O‘";case "ğ":return "g‘";case "Ğ":return "G‘";case "ş":return "sh";case "Ş":return caps?"SH":"Sh";case "ç":return "ch";default:return caps?"CH":"Ch";}
            });}
            var b=new StringBuilder();
            for(int i=0;i<s.Length;i++){
                char c=s[i],low=char.ToLowerInvariant(c),next=i+1<s.Length?s[i+1]:'\0';bool upper=char.IsUpper(c);char combined='\0';
                if((low=='o'||low=='g')&&next!='\0'&&Apostrophes.IndexOf(next)>=0&&(Letter(s,i+2)||(low=='g'&&Letter(s,i-1)))) combined=low=='o'?'ö':'ğ';
                else if(char.ToLowerInvariant(next)=='h'&&(low=='s'||low=='c'))combined=low=='s'?'ş':'ç';
                if(combined!='\0'){b.Append(upper?char.ToUpperInvariant(combined):combined);i++;}
                else if(c!='ʼ'&&Apostrophes.IndexOf(c)>=0&&Letter(s,i-1)&&Letter(s,i+1))b.Append('ʼ');else b.Append(c);
            }return b.ToString();
        });
    }
    public static void Test(){
        string[,] cases={{"O‘zbekiston, g‘alaba, shahar, chiroyli.","Özbekiston, ğalaba, şahar, çiroyli."},{"O'ZBEK SHAHAR CHIROQ TOG'","ÖZBEK ŞAHAR ÇIROQ TOĞ"},{"ma'no san’at Isʼhoq","maʼno sanʼat Isʼhoq"},{"https://example.com/sh user@shahar.uz shahar","https://example.com/sh user@shahar.uz şahar"},{"'o' va 'g'","'o' va 'g'"}};
        for(int i=0;i<cases.GetLength(0);i++)if(Convert(cases[i,0],false)!=cases[i,1])throw new Exception("Converter self-test "+i);
        if(Convert("Özbekiston, ŞAHAR, Çiroq",true)!="O‘zbekiston, SHAHAR, Chiroq")throw new Exception("Reverse converter test");
    }
}

sealed class HarfWindow:Form {
    readonly NotifyIcon tray; readonly CheckBox enabled,startup; readonly Label state;
    readonly TextBox source,result; bool closing,dispatching;
    IntPtr hook,mouseHook,lastWindow; Native.HookProc keyProc,mouseProc;
    char previous; int lastTick;
    const int ReplaceMessage=0x8001; char pending;
    public HarfWindow(){
        Text="Harf — Windows yordamchisi (Preview)";Width=640;Height=610;MinimumSize=new Size(640,610);StartPosition=FormStartPosition.CenterScreen;Font=new Font("Segoe UI",10);BackColor=Color.FromArgb(248,248,242);
        var header=new Label{Text="Harf.  Yangi alifbo, tanish klaviatura.",Left=22,Top=20,Width=580,Height=40,Font=new Font("Segoe UI",18,FontStyle.Bold)};Controls.Add(header);
        state=new Label{Left=24,Top=70,Width=570,Height=42};Controls.Add(state);
        enabled=new CheckBox{Text="Avtoalmashtirish (US lotin klaviaturasi uchun)",Left=24,Top=115,Width=540};enabled.CheckedChanged+=(s,e)=>{previous='\0';UpdateState();};Controls.Add(enabled);
        startup=new CheckBox{Text="Windows bilan ishga tushirish",Left=24,Top=148,Width=540};
        using(var run=Registry.CurrentUser.OpenSubKey(@"Software\Microsoft\Windows\CurrentVersion\Run"))startup.Checked=run!=null&&run.GetValue("Harf")!=null;
        startup.CheckedChanged+=(s,e)=>{try{using(var run=Registry.CurrentUser.CreateSubKey(@"Software\Microsoft\Windows\CurrentVersion\Run")){if(startup.Checked)run.SetValue("Harf","\""+Application.ExecutablePath+"\"");else run.DeleteValue("Harf",false);}}catch(Exception ex){MessageBox.Show(ex.Message,"Harf");}};Controls.Add(startup);
        Controls.Add(new Label{Text="AltGr + O/G/S/C → ö/ğ/ş/ç • Shift bilan katta harf\nCtrl + Alt + Y → avtoalmashtirishni yoqish/o‘chirish",Left=24,Top=185,Width=570,Height=48});
        Controls.Add(new Label{Text="Matn konvertori",Left=24,Top=245,Width=550});
        source=new TextBox{Left=24,Top=273,Width=570,Height=85,Multiline=true,ScrollBars=ScrollBars.Vertical};Controls.Add(source);
        var forward=new Button{Text="Eski → Yangi",Left=24,Top=369,Width=140,Height=32};forward.Click+=(s,e)=>result.Text=Converter.Convert(source.Text,false);Controls.Add(forward);
        var reverse=new Button{Text="Yangi → Eski",Left=176,Top=369,Width=140,Height=32};reverse.Click+=(s,e)=>result.Text=Converter.Convert(source.Text,true);Controls.Add(reverse);
        result=new TextBox{Left=24,Top=412,Width=570,Height=85,Multiline=true,ReadOnly=true,ScrollBars=ScrollBars.Vertical};Controls.Add(result);
        var copy=new Button{Text="Nusxalash",Left=24,Top=510,Width=140,Height=32};copy.Click+=(s,e)=>{if(result.TextLength>0)try{Clipboard.SetText(result.Text);}catch{MessageBox.Show("Clipboard band. Qayta urinib ko‘ring.","Harf");}};Controls.Add(copy);
        var exit=new Button{Text="Dasturdan chiqish",Left=430,Top=510,Width=164,Height=32};exit.Click+=(s,e)=>{closing=true;Close();};Controls.Add(exit);
        var menu=new ContextMenuStrip();menu.Items.Add("Oynani ochish",null,(s,e)=>{Show();WindowState=FormWindowState.Normal;Activate();});menu.Items.Add("Avtoalmashtirishni almashtirish",null,(s,e)=>enabled.Checked=!enabled.Checked);menu.Items.Add("Chiqish",null,(s,e)=>{closing=true;Close();});
        tray=new NotifyIcon{Text="Harf — avtomatik rejim o‘chiq",Icon=SystemIcons.Information,ContextMenuStrip=menu,Visible=true};tray.DoubleClick+=(s,e)=>{Show();Activate();};
        // Force Handle creation before installing callbacks / registering hotkeys.
        IntPtr handle=Handle;
        bool hotkeys=true;
        for(int i=0;i<4;i++){int key="OGSC"[i];hotkeys&=Native.RegisterHotKey(handle,10+i,0x4003,(uint)key);hotkeys&=Native.RegisterHotKey(handle,20+i,0x4007,(uint)key);}
        hotkeys&=Native.RegisterHotKey(handle,30,0x4003,0x59);
        keyProc=Keyboard;mouseProc=Mouse;
        using(var process=Process.GetCurrentProcess())using(var module=process.MainModule){IntPtr instance=Native.GetModuleHandle(module.ModuleName);hook=Native.SetWindowsHookEx(13,keyProc,instance,0);mouseHook=Native.SetWindowsHookEx(14,mouseProc,instance,0);}
        if(!hotkeys||hook==IntPtr.Zero||mouseHook==IntPtr.Zero)MessageBox.Show("Some keyboard shortcuts or hooks are unavailable. Close conflicting keyboard tools and restart Harf. Text conversion remains available.","Harf");
        UpdateState();
    }
    void UpdateState(){if(state==null)return;state.Text=enabled.Checked?"Avtoalmashtirish YOQILGAN. Parol va maxfiy maydonlarda o‘chirib turing.":"Avtoalmashtirish O‘CHIQ. AltGr tugmalari ishlaydi.";if(tray!=null)tray.Text=enabled.Checked?"Harf — avtomatik rejim yoqilgan":"Harf — avtomatik rejim o‘chiq";}
    bool Blocked(IntPtr window){
        try{uint pid;uint thread=Native.GetWindowThreadProcessId(window,out pid);string name=Process.GetProcessById((int)pid).ProcessName.ToLowerInvariant();
            if(Regex.IsMatch(name,@"^(code|code-insiders|devenv|idea64|pycharm64|webstorm64|goland64|rider64|notepad\+\+|windowsterminal|powershell|pwsh|cmd|mintty|putty)$"))return true;
            var info=new Native.GUIINFO();info.size=Marshal.SizeOf(typeof(Native.GUIINFO));
            if(Native.GetGUIThreadInfo(thread,ref info)&&info.focus!=IntPtr.Zero){if(Native.SendMessage(info.focus,0xD2,IntPtr.Zero,IntPtr.Zero)!=IntPtr.Zero)return true;}
        }catch{return true;}return false;
    }
    IntPtr Keyboard(int code,IntPtr wp,IntPtr lp){
        if(code>=0&&!dispatching){var k=(Native.KEYDATA)Marshal.PtrToStructure(lp,typeof(Native.KEYDATA));
            if((k.flags&0x10)==0&&(wp.ToInt64()==0x100||wp.ToInt64()==0x104)){
                IntPtr window=Native.GetForegroundWindow();int now=Environment.TickCount;
                if(window!=lastWindow||unchecked(now-lastTick)>2500)previous='\0';lastWindow=window;lastTick=now;
                if(!enabled.Checked||Blocked(window)||Down(0x11)||Down(0x12)||Down(0x5B)||Down(0x5C))previous='\0';
                else {char current='\0';if(k.key>=65&&k.key<=90){bool upper=Down(0x10)^((Native.GetKeyState(0x14)&1)!=0);current=(char)(upper?k.key:k.key+32);}else if(k.key==0xDE&&!Down(0x10))current='\'';else if(k.key==0xC0&&!Down(0x10))current='`';
                    char lower=char.ToLowerInvariant(previous),lowerCurrent=char.ToLowerInvariant(current),combined='\0';
                    if(lowerCurrent=='h'&&(lower=='s'||lower=='c'))combined=lower=='s'?'ş':'ç';
                    else if((current=='\''||current=='`')&&(lower=='o'||lower=='g'))combined=lower=='o'?'ö':'ğ';
                    if(combined!='\0'){pending=char.IsUpper(previous)?char.ToUpperInvariant(combined):combined;previous='\0';Native.PostMessage(Handle,ReplaceMessage,IntPtr.Zero,IntPtr.Zero);return new IntPtr(1);}previous=current;
                }
            }
        }return Native.CallNextHookEx(hook,code,wp,lp);
    }
    IntPtr Mouse(int code,IntPtr wp,IntPtr lp){if(code>=0&&(wp.ToInt64()==0x201||wp.ToInt64()==0x204||wp.ToInt64()==0x207))previous='\0';return Native.CallNextHookEx(mouseHook,code,wp,lp);}
    static bool Down(int key){return (Native.GetAsyncKeyState(key)&0x8000)!=0;}
    protected override void WndProc(ref Message message){
        if(message.Msg==0x312){int id=message.WParam.ToInt32();if(id==30)enabled.Checked=!enabled.Checked;else if(id>=10&&id<=13||id>=20&&id<=23){char c="öğşç"[id%10];if(id>=20)c=char.ToUpperInvariant(c);Send(c,false);}previous='\0';}
        else if(message.Msg==ReplaceMessage){if(Native.GetForegroundWindow()==lastWindow)Send(pending,true);}
        base.WndProc(ref message);
    }
    void Send(char c,bool backspace){
        // Wait for shortcut modifiers to be physically released before Unicode injection.
        if(!backspace&&(Down(0x11)||Down(0x12)||Down(0x10))){IntPtr target=Native.GetForegroundWindow();var timer=new System.Windows.Forms.Timer{Interval=15};int ticks=0;timer.Tick+=(s,e)=>{if(Native.GetForegroundWindow()!=target){timer.Stop();timer.Dispose();}else if(!Down(0x11)&&!Down(0x12)&&!Down(0x10)){timer.Stop();timer.Dispose();Send(c,false);}else if(++ticks>100){timer.Stop();timer.Dispose();}};timer.Start();return;}
        dispatching=true;try{var keys=new Native.INPUT[backspace?4:2];int n=0;if(backspace){keys[n++]=Native.Input(8,0,0);keys[n++]=Native.Input(8,0,2);}keys[n++]=Native.Input(0,c,4);keys[n++]=Native.Input(0,c,6);if(Native.SendInput((uint)keys.Length,keys,Marshal.SizeOf(typeof(Native.INPUT)))!=keys.Length)tray.ShowBalloonTip(2000,"Harf","This app may not allow input injection. Try the text converter.",ToolTipIcon.Info);}finally{dispatching=false;}
    }
    protected override void OnFormClosing(FormClosingEventArgs e){if(!closing&&e.CloseReason==CloseReason.UserClosing){e.Cancel=true;Hide();return;}Native.UnhookWindowsHookEx(hook);Native.UnhookWindowsHookEx(mouseHook);for(int i=10;i<=30;i++)Native.UnregisterHotKey(Handle,i);tray.Visible=false;tray.Dispose();base.OnFormClosing(e);}
}

static class Native {
    public delegate IntPtr HookProc(int n,IntPtr w,IntPtr l);
    [StructLayout(LayoutKind.Sequential)]public struct KEYDATA{public uint key,scan,flags,time;public UIntPtr extra;}
    [StructLayout(LayoutKind.Sequential)]public struct KEYINPUT{public ushort key,scan;public uint flags,time;public UIntPtr extra;}
    [StructLayout(LayoutKind.Sequential)]public struct MOUSEINPUT{public int x,y;public uint data,flags,time;public UIntPtr extra;}
    [StructLayout(LayoutKind.Explicit)]public struct UNION{[FieldOffset(0)]public KEYINPUT keyboard;[FieldOffset(0)]public MOUSEINPUT mouse;}
    [StructLayout(LayoutKind.Sequential)]public struct INPUT{public uint type;public UNION data;}
    [StructLayout(LayoutKind.Sequential)]public struct RECT{public int left,top,right,bottom;}
    [StructLayout(LayoutKind.Sequential)]public struct GUIINFO{public int size;public uint flags;public IntPtr active,focus,capture,menu,move,caret;public RECT rect;}
    public static INPUT Input(ushort key,char scan,uint flags){return new INPUT{type=1,data=new UNION{keyboard=new KEYINPUT{key=key,scan=scan,flags=flags}}};}
    public static INPUT Input(ushort key,int scan,uint flags){return Input(key,(char)scan,flags);}
    [DllImport("user32.dll",SetLastError=true)]public static extern IntPtr SetWindowsHookEx(int id,HookProc p,IntPtr module,uint thread);
    [DllImport("user32.dll")]public static extern bool UnhookWindowsHookEx(IntPtr hook);
    [DllImport("user32.dll")]public static extern IntPtr CallNextHookEx(IntPtr hook,int code,IntPtr w,IntPtr l);
    [DllImport("user32.dll")]public static extern short GetAsyncKeyState(int key);
    [DllImport("user32.dll")]public static extern short GetKeyState(int key);
    [DllImport("user32.dll")]public static extern IntPtr GetForegroundWindow();
    [DllImport("user32.dll")]public static extern uint GetWindowThreadProcessId(IntPtr w,out uint pid);
    [DllImport("user32.dll")]public static extern bool GetGUIThreadInfo(uint thread,ref GUIINFO info);
    [DllImport("user32.dll")]public static extern IntPtr SendMessage(IntPtr w,int msg,IntPtr a,IntPtr b);
    [DllImport("user32.dll")]public static extern bool PostMessage(IntPtr w,int msg,IntPtr a,IntPtr b);
    [DllImport("kernel32.dll",CharSet=CharSet.Unicode)]public static extern IntPtr GetModuleHandle(string name);
    [DllImport("user32.dll",SetLastError=true)]public static extern uint SendInput(uint count,INPUT[] input,int size);
    [DllImport("user32.dll",SetLastError=true)]public static extern bool RegisterHotKey(IntPtr w,int id,uint modifiers,uint key);
    [DllImport("user32.dll")]public static extern bool UnregisterHotKey(IntPtr w,int id);
}
